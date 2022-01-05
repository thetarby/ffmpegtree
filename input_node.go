package ffmpegtree

import (
	"fmt"
	"time"
)

type IInputNode interface {
	INode
	ToString() string
	GetInputIdx() int
	SetInputIdx(int)
}

type InputNode struct {
	BaseNode
	InputName   string
	Offset, Len *time.Duration
	inputIdx    int
}

func (i *InputNode) ToString() string {
	str := fmt.Sprintf("-i %v", i.InputName)
	if i.Offset != nil {
		str = fmt.Sprintf("-ss %v %v", fmtDuration(*i.Offset), str)
	}
	if i.Len != nil {
		str = fmt.Sprintf("-t %v %v", fmtDuration(*i.Len), str)
	}

	return str
}

func (i *InputNode) GetInputIdx() int {
	return i.inputIdx
}

func (i *InputNode) SetInputIdx(idx int) {
	i.inputIdx = idx
}

func NewInputNode(name string, len, offset *time.Duration) *InputNode {
	return &InputNode{
		BaseNode:  NewBaseNode(nil),
		InputName: name,
		Offset:    offset,
		Len:       len,
	}
}

// ISelectStreamNode is still an IInputNode, but it is combined with another IInputNode to select a specific stream
// from the input
type ISelectStreamNode interface {
	INode
	GetStream() string
	GetInputNode() IInputNode
}

// SelectStreamNode implements ISelectStreamNode
var _ ISelectStreamNode = &SelectStreamNode{}

type SelectStreamNode struct {
	BaseNode
	input IInputNode
	idx   int
}

func (s *SelectStreamNode) GetInputNode() IInputNode {
	return s.input
}

func (s *SelectStreamNode) GetStream() string {
	return fmt.Sprintf("[%v:%v]", s.input.GetInputIdx(), s.idx)
}

func NewSelectStreamNode(input IInputNode, idx int) *SelectStreamNode {
	return &SelectStreamNode{
		BaseNode: NewBaseNode([]INode{input}),
		input:    input,
		idx:      idx,
	}
}
