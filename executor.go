package ffmpegtree

import (
	"fmt"
	"strings"
)

// Chain is a singly linked chain of IFilterNode's which is written seperated by ',' in ffmpeg commands
type Chain []IFilterNode

// ToString returns string representation of the chain of filters seperated by ','
func (c Chain) ToString(includeOutStream bool) string {
	filterStrings := make([]string, 0)
	for i := len(c) - 1; i >= 0; i-- {
		filterStrings = append(filterStrings, FilterNodeToStr(c[i]))
	}

	// if ends with split
	res := ""
	if split, ok := c[0].(*SplitNode); ok {
		outs := ""
		for i := 0; i < split.fanOut; i++ {
			outs += fmt.Sprintf("[%v]", split.GetOutStreamName())
		}
		res = strings.Join(filterStrings, ",") + outs
	} else {
		res = strings.Join(filterStrings, ",")
		if includeOutStream {
			res += fmt.Sprintf("[%v]", c[0].GetOutStreamName())
		}
	}

	return res
}

type FFmpegExecutor struct {
	acc        []string
	inputs     []IInputNode
	visited    map[string]bool
	dependents *DependentsMap
	q          []INode
	maps       []IMap
	outName    string
	outOptions []string
}

func (e *FFmpegExecutor) ToFfmpeg(nodes ...INode) FfmpegCommad {
	for _, node := range nodes {
		// find all IInputNode and assign their indexes. it will affect the order they show up in the output
		e.setInputIdx(node)

		// first preprocess tree
		// insert select stream nodes where it is missing
		e.insertSelectStream(node)

		// insert split nodes if a stream is input to more than one node
		e.insertSplit(node)
	}

	// each node can access its inputs but cannot access to nodes which depends on itself. traverse tree
	// and save dependencies in a hashmap structure. it will be useful while executing tree.
	e.dependents = GetDependents(nodes...)

	// start traversal from root node
	e.q = nodes
	r := e.toFfmpeg()

	// generate input options which are in the form of "-i ***.mp4"
	inputs := make([]string, 0, len(e.inputs))
	for _, input := range e.inputs {
		inputs = append(inputs, input.ToString()...)
	}

	// generate map options which are in the form of "-map '0:0'" or "-map '[var_1:0]'"
	maps := make([]string, 0, len(e.maps))
	for _, iMap := range e.maps {
		maps = append(maps, iMap.ToString()...)
	}

	// put it all together
	res := make([]string, 0)
	res = append(res, inputs...)
	res = append(res, "-filter_complex", r)
	res = append(res, maps...)
	res = append(res, e.outOptions...)
	res = append(res, e.outName)
	return res
}

func (e *FFmpegExecutor) toFfmpeg() string {
	// NOTE: we are traversing tree in bfs manner and as lons as each node only has 1 dependents, it normally guarentees that a filter or a stream variable
	// shows up before all of its dependents in the output, avoiding 'forward declarations, which ffmpeg script does not support'.
	// It is not guaranteed for split nodes because they could have more than one dependents. A split node may have a child node
	// which still has a dependent unprocessed node. To guarantee that the node which is splitted comes before all of its
	// dependents, we need to process its split node after all of its dependents are queued.

	// nodeEncounters is a mapping from nodes to their fan out number(dependent count). Each time a node is popped
	// up from queue, its corresponding mapping is decreased by one and when it hits 0, it means the last dependent
	// is processed hence we can proceed to process the current node and add it to output
	nodeEncounters := make(map[string]int)

	for len(e.q) > 0 {
		tree := e.q[0]
		e.q = e.q[1:]
		if e.visited[tree.GetID()] {
			continue
		}

		switch tree.(type) {
		case IFilterNode:
			if _, ok := nodeEncounters[tree.GetID()]; !ok {
				nodeEncounters[tree.GetID()] = len(e.dependents.Get(tree))
			}

			if nodeEncounters[tree.GetID()] = nodeEncounters[tree.GetID()] - 1; nodeEncounters[tree.GetID()] > 0 {
				// if not 0 delay processing of the node since it migth still have unprocessed dependents
				continue
			}

			// every IFilterNode can be treated as a chain of IFilterNode's even if it consists of only one node
			c, newTree := e.toChain(tree)
			e.visited[tree.GetID()] = true
			tree = newTree

			f := ""
			for _, node := range tree.GetInputs() {
				parentFilterNode, ok := node.(IFilterNode)
				if !ok {
					ssn := node.(ISelectStreamNode)
					f += ssn.GetOutStreamName()
					continue
				}
				f += fmt.Sprintf("[%v]", parentFilterNode.GetOutStreamName())
			}

			f += c.ToString(len(e.dependents.Get(c[0])) > 0 || e.isMapped(c[0]))
			e.acc = append([]string{f}, e.acc...)

		case IInputNode:
			e.visited[tree.GetID()] = true
		}

		e.q = append(e.q, tree.GetInputs()...)
	}

	return strings.Join(e.acc, ";")
}

// toChain creates a chain of IFilterNode's by adding given node then, starts to follow every node's inputs and adds it to chain
// if the node has only one child and one dependency
func (e *FFmpegExecutor) toChain(tree INode) (c Chain, ret INode) {
	stack := make(Chain, 0)
	stack = append(stack, tree.(IFilterNode))
	for ; len(tree.GetInputs()) == 1; tree = tree.GetInputs()[0] {
		curr := tree.GetInputs()[0]
		if e.visited[curr.GetID()] || len(e.dependents.Get(curr)) > 1 {
			break
		}
		fn, ok := curr.(IFilterNode)
		if !ok {
			break
		}
		stack = append(stack, fn)
	}

	return stack, stack[len(stack)-1]
}

// insertSelectStream traverses the graph and inserts an ISelectStreamNode when an IInputNode is directly fed into a
// IFilterNode. Default selected stream is '0'
func (e *FFmpegExecutor) insertSelectStream(t INode) {
	d := GetDependents(t)
	for _, s := range d.Keys() {
		nodes := d.Get(s)
		in, ok := s.(IInputNode)
		if !ok {
			continue
		}
		ssn := NewSelectStreamNode(in, 0)
		for _, node := range nodes {
			_, ok := node.(IFilterNode)
			if ok {
				// if input is directly connected to an IFilterNode insert ssn in between
				inps := node.GetInputs()
				i := 0
				for ; i < len(inps); i++ {
					if inps[i].GetID() == s.GetID() {
						break
					}
				}
				inps[i] = ssn
				node.SetInputs(inps)
			}
		}
	}
}

// insertSplit traverses the graph and inserts an ISplitNode if a stream is used as input for more than one times.
// This is required by ffmpeg syntax.
func (e *FFmpegExecutor) insertSplit(t INode) {
	d := GetDependents(t)
	for _, currNode := range d.Keys() {
		if _, ok := currNode.(ISplitNode); ok {
			continue
		}

		dependents := d.Get(currNode)
		_, ok := currNode.(IInputNode) // input node can be input to more than once without splitting
		if len(dependents) > 1 && !ok {
			split := NewSplitNode(currNode, len(dependents))
			for _, node := range dependents {
				inps := node.GetInputs()
				i := 0
				for ; i < len(inps); i++ {
					if inps[i].GetID() == currNode.GetID() {
						break
					}
				}
				inps[i] = split
				node.SetInputs(inps)
			}
		}
	}
}

// setInputIdx traverses the graph from a given node and discovers all input nodes and assign them an index.
func (e *FFmpegExecutor) setInputIdx(t INode) {
	d := GetDependents(t)
	for _, s := range d.Keys() {
		in, ok := s.(IInputNode)
		if ok && !e.isInInputs(in) {
			in.SetInputIdx(len(e.inputs))
			e.inputs = append(e.inputs, in)
		}
	}

	// add input nodes which are mapped but is not used in graph at all
	for _, iMap := range e.maps {
		n := iMap.GetStreamNode()
		in, ok := n.(IInputNode)
		if ok && d.Get(in) == nil && !e.isInInputs(in) {
			in.SetInputIdx(len(e.inputs))
			e.inputs = append(e.inputs, in)
		}
	}
}

func (e *FFmpegExecutor) isMapped(n INode) bool {
	for _, m := range e.maps {
		if m.GetStreamNode().GetID() == n.GetID() {
			return true
		}
	}

	return false
}

func (e *FFmpegExecutor) isInInputs(n IInputNode) bool {
	for _, inp := range e.inputs {
		if inp.GetID() == n.GetID() {
			return true
		}
	}

	return false
}

func NewFfmpegExecutor(maps []IMap, outName string, outOptions []string) *FFmpegExecutor {
	return &FFmpegExecutor{
		acc:        nil,
		inputs:     nil,
		visited:    make(map[string]bool),
		dependents: NewDependentsMap(),
		q:          nil,
		maps:       maps,
		outName:    outName,
		outOptions: outOptions,
	}
}

type FfmpegCommad []string

func (cmd *FfmpegCommad) FilterComplex() string {
	for i, arg := range *cmd {
		if strings.Contains(arg, "filter_complex") {
			return []string(*cmd)[i+1]
		}
	}

	return ""
}
