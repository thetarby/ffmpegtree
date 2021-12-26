package ffmpegtree

import "fmt"

type ISplitNode interface {
	INode
	IFilterNode
	GetFanOut() int
}

// SplitNode implements ISplitNode
var _ ISplitNode = &SplitNode{}

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
	return "split"
}

func NewSplitNode(input INode, fanOut int) *SplitNode {
	return &SplitNode{
		BaseFilterNode: *NewBaseFilterNode([]INode{input}, randStr()),
		fanOut:         fanOut,
	}
}
