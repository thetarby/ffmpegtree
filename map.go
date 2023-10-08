package ffmpegtree

import (
	"fmt"
)

type IMap interface {
	ToString() []string

	// GetStreamNode returns the node from which stream is mapped
	GetStreamNode() INode
}

var _ IMap = &MapFromFilterNode{}

type MapFromFilterNode struct {
	filterNode IFilterNode
}

func (m *MapFromFilterNode) GetStreamNode() INode {
	return m.filterNode
}

func (m *MapFromFilterNode) ToString() []string {
	return []string{"-map", fmt.Sprintf("[%v]", m.filterNode.GetOutStreamName())}
}

var _ IMap = &MapFromInputNode{}

type MapFromInputNode struct {
	input  IInputNode
	stream string
}

func (m *MapFromInputNode) GetStreamNode() INode {
	return m.input
}

func (m *MapFromInputNode) ToString() []string {
	return []string{"-map", fmt.Sprintf("%v:%v", m.input.GetInputIdx(), m.stream)}
}

func NewMap(fromNode INode, opts ...string) IMap {
	in, ok := fromNode.(IInputNode)
	if ok {
		return &MapFromInputNode{
			input:  in,
			stream: opts[0],
		}
	}

	fn, ok := fromNode.(IFilterNode)
	if ok {
		return &MapFromFilterNode{
			filterNode: fn,
		}
	}

	return nil
}

func Select(nodes []INode, outName string, outputOptions []string, maps ...IMap) FfmpegCommand {
	exec := NewFfmpegExecutor(maps, outName, outputOptions)
	res := exec.ToFfmpeg(nodes...)
	return res
}
