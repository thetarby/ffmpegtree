package ffmpegtree

import "fmt"

type TimelineAcceptingFilterNode struct {
	BaseFilterNode
	enableExpr   string
	since, until *float64
}

func (n *TimelineAcceptingFilterNode) Enable(exp string) {
	n.enableExpr = exp
}

func (n *TimelineAcceptingFilterNode) Since(since float64) {
	n.since = &since
}

func (n *TimelineAcceptingFilterNode) Until(until float64) {
	n.until = &until
}

func (n *TimelineAcceptingFilterNode) EnableExpr() string {
	if n.enableExpr != "" {
		return n.enableExpr
	}

	if n.enableExpr == "" && n.since == nil && n.until == nil {
		return ""
	}

	if n.since != nil && n.until != nil {
		return fmt.Sprintf("between(t, %.2f, %.2f)", *n.since, *n.until)
	}

	if n.since != nil {
		return fmt.Sprintf("gte(t, %.2f)", *n.since)
	}

	// if only until is set
	return fmt.Sprintf("lte(t, %.2f)", *n.until)
}

func NewTimelineAcceptingFilterNode(children []INode, outStreamName string) *TimelineAcceptingFilterNode {
	return &TimelineAcceptingFilterNode{
		BaseFilterNode: *NewBaseFilterNode(children, outStreamName),
	}
}
