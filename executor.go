package ffmpegtree

import (
	"fmt"
	"strings"
)

// Chain is a singly linked chain of IFilterNodes which is written seperated by ',' in ffmpeg commands
type Chain []IFilterNode

// ToString returns string representation of the chain of filters seperated by ','
func (c Chain) ToString(includeOutStream bool) string {
	filterStrings := make([]string, 0)
	for i := len(c) - 1; i >= 0; i-- {
		filterStrings = append(filterStrings, c[i].FilterString())
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
}

func (e *FFmpegExecutor) ToFfmpeg(tree INode) string {
	// first preprocess tree
	// insert select stream nodes where it is missing
	e.insertSelectStream(tree)

	// insert split nodes if a stream is input to more than one node
	e.insertSplit(tree)

	// find all IInputNode and assign their indexes. it will affect the order they show up in the output
	e.setInputIdx(tree)

	// each node can access its inputs but cannot access to nodes which depends on itself. traverse tree
	// and save dependencies in a hashmap structure. it will be useful while executing tree.
	e.dependents = GetDependents(tree)

	// start traversal from root node
	e.q = []INode{tree}
	r := e.toFfmpeg()

	// generate input options which are in the form of "-i ***.mp4"
	inputs := ""
	for _, input := range e.inputs {
		inputs += input.ToString() + " "
	}

	// generate map options which are in the form of "-map '0:0'" or "-map '[var_1:0]'"
	maps := ""
	for _, iMap := range e.maps {
		maps += fmt.Sprintf(" %v", iMap.ToString())
	}

	// put it all together
	return inputs + " -filter_complex '" + r + "'" + maps + fmt.Sprintf(" %v", e.outName)
}

func (e *FFmpegExecutor) toFfmpeg() string {
	for len(e.q) > 0 {
		tree := e.q[0]
		e.q = e.q[1:]
		if e.visited[tree.GetID()] {
			continue
		}

		switch tree.(type) {
		case IFilterNode:
			e.visited[tree.GetID()] = true

			// every IFilterNode can be treated as a chain of IFilterNode's even if it consists of only one node
			c, newTree := e.toChain(tree)
			tree = newTree

			if len(tree.GetInputs()) > 0 { // TODO: unnecessary if
				f := ""
				for _, node := range tree.GetInputs() {
					parentFilterNode, ok := node.(IFilterNode)
					if !ok {
						ssn := node.(ISelectStreamNode) // TODO: check if it is input node
						f += ssn.GetStream()            // TODO: should be able to select any stream from input
						continue
					}
					f += fmt.Sprintf("[%v]", parentFilterNode.GetOutStreamName())
				}

				f += c.ToString(len(e.dependents.Get(c[0])) > 0 || e.isMapped(c[0]))
				e.acc = append([]string{f}, e.acc...)
			}
		case IInputNode:
			e.visited[tree.GetID()] = true
		}

		e.q = append(e.q, tree.GetInputs()...)
	}

	return strings.Join(e.acc, ";")
}

// toChain creates a chain by adding given node then, starts to follow every node's inputs and adds it to chain
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
		nodes, _ := d.GetAndCheck(s)
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
	for _, s := range d.Keys() {
		nodes, _ := d.GetAndCheck(s)
		_, ok := s.(IInputNode)
		if len(nodes) > 1 && !ok {
			split := NewSplitNode(s, len(nodes))
			for _, node := range nodes {
				inps := node.GetInputs()
				i := 0
				for ; i < len(inps); i++ {
					if inps[i].GetID() == s.GetID() {
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
		if ok {
			in.SetInputIdx(len(e.inputs))
			e.inputs = append(e.inputs, in)
		}
	}

	// add input nodes which are mapped but is not used in graph at all
	for _, iMap := range e.maps {
		n := iMap.GetStreamNode()
		in, ok := n.(IInputNode)
		if ok && d.Get(in) == nil {
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

func NewFfmpegExecutor(maps []IMap, outName string) *FFmpegExecutor {
	return &FFmpegExecutor{
		acc:        nil,
		inputs:     nil,
		visited:    make(map[string]bool),
		dependents: NewDependentsMap(),
		q:          nil,
		maps:       maps,
		outName:    outName,
	}
}
