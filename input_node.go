package ffmpegtree

import (
	"fmt"
	"strconv"
	"time"
)

type IInputNode interface {
	INode
	ToString() []string
	GetInputIdx() int
	SetInputIdx(int)
	GetInputName() string
}

// InputNode implements IInputNode
var _ IInputNode = &InputNode{}

type InputNode struct {
	BaseNode
	InputName   string
	Offset, Len *time.Duration
	inputIdx    int
	isLoop      bool
}

func (i *InputNode) ToString() []string {
	res := make([]string, 0)
	res = append(res, "-i", i.InputName)
	if i.Offset != nil {
		res = append([]string{"-ss", fmtDuration(*i.Offset)}, res...)
	}
	if i.Len != nil {
		res = append([]string{"-t", fmtDuration(*i.Len)}, res...)
	}
	if i.isLoop {
		res = append([]string{"-stream_loop", "-1"}, res...)
	}

	return res
}

func (i *InputNode) GetInputIdx() int {
	return i.inputIdx
}

func (i *InputNode) SetInputIdx(idx int) {
	i.inputIdx = idx
}

func (i *InputNode) GetInputName() string {
	return i.InputName
}

func NewInputNode(name string, len, offset *time.Duration) *InputNode {
	return &InputNode{
		BaseNode:  NewBaseNode(nil),
		InputName: name,
		Offset:    offset,
		Len:       len,
	}
}

// NewAudioInputNode opens an input file and selects audio stream.
func NewAudioInputNode(name string, len, offset *time.Duration) INode {
	return NewSelectStreamNode(NewInputNode(name, len, offset), AudioStream)
}

func NewInputNodeLoop(name string) *InputNode {
	return &InputNode{
		BaseNode:  NewBaseNode(nil),
		InputName: name,
		isLoop:    true,
	}
}

const VideoStream = -1
const AudioStream = -2

// ISelectStreamNode is still an IInputNode, but it is combined with another IInputNode to select a specific stream
// from the input
type ISelectStreamNode interface {
	INode
	Streamer
	GetInputNode() IInputNode
}

// SelectStreamNode implements ISelectStreamNode
var _ ISelectStreamNode = &SelectStreamNode{}

type SelectStreamNode struct {
	BaseNode
	input IInputNode
	idx   string
}

func (s *SelectStreamNode) GetInputNode() IInputNode {
	return s.input
}

func (s *SelectStreamNode) GetOutStreamName() string {
	return fmt.Sprintf("[%v:%v]", s.input.GetInputIdx(), s.idx)
}

func NewSelectStreamNode(input IInputNode, idx int) *SelectStreamNode {
	str := strconv.Itoa(idx)
	if idx == AudioStream {
		str = "a"
	} else if idx == VideoStream {
		str = "v"
	}

	return &SelectStreamNode{
		BaseNode: NewBaseNode([]INode{input}),
		input:    input,
		idx:      str,
	}
}
