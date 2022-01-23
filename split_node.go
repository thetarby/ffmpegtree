package ffmpegtree

import "fmt"

type ISplitNode interface {
	INode
	IFilterNode
	GetFanOut() int
}

// SplitNode implements ISplitNode
var _ ISplitNode = &SplitNode{}

// NOTE: to split an audio stream asplit filter should be used but for now SplitNode only works for
// video streams since there is no audio-video stream separation as of now. So splitting audio streams
// would not work.
type SplitNode struct {
	BaseFilterNode
	fanOut, called int
}

func (s *SplitNode) GetFanOut() int {
	return s.fanOut
}

func (s *SplitNode) GetOutStreamName() string {
	defer func() { s.called++ }()
	return fmt.Sprintf("%v_%v", s.OutStreamName, s.called%s.fanOut)
}

func (b *SplitNode) FilterString() string {
	if b.fanOut > 2{
		return fmt.Sprintf("split=%v", b.fanOut)
	}
	return "split"
}

func NewSplitNode(input INode, fanOut int) *SplitNode {
	return &SplitNode{
		BaseFilterNode: *NewBaseFilterNode([]INode{input}, randStr()),
		fanOut:         fanOut,
	}
}
