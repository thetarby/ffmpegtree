package ffmpegtree

import "fmt"

type IMergeNode interface {
	INode
	IFilterNode
	GetFanIn() int
}

// AmergeNode implements IMergeNode
var _ IMergeNode = &AmergeNode{}

type AmergeNode struct {
	BaseFilterNode
	inpCount, called int
}

func (s *AmergeNode) GetFanIn() int {
	return s.inpCount
}

func (s *AmergeNode) GetOutStreamName() string {
	defer func() { s.called++ }()
	return fmt.Sprintf("%v_%v", s.OutStreamName, s.called%s.inpCount)
}

func (s *AmergeNode) FilterString() string {
	return fmt.Sprintf("amerge=inputs=%v", len(s.BaseFilterNode.inputs))
}

func NewMergeNode(inputs ...INode) *AmergeNode {
	return &AmergeNode{
		BaseFilterNode: *NewBaseFilterNode(inputs, randStr()),
	}
}
