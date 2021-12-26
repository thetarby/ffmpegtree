package ffmpegtree

import "strconv"

type INode interface {
	GetID() string
	GetInputs() []INode
	SetInputs([]INode)
}

type BaseNode struct {
	id     string
	inputs []INode
}

func (b *BaseNode) GetID() string {
	return b.id
}

func (b *BaseNode) GetInputs() []INode {
	return b.inputs
}

func (b *BaseNode) SetInputs(children []INode) {
	b.inputs = children
}

var id = 0

func NewBaseNode(children []INode) BaseNode {
	id++
	return BaseNode{
		id:     strconv.Itoa(id),
		inputs: children,
	}
}

func GetDepending(t INode) map[INode][]INode {
	res := make(map[INode][]INode)
	visited := make(map[string]bool)
	getDepending(t, res, visited)
	return res
}

func getDepending(t INode, acc map[INode][]INode, visited map[string]bool) {
	if visited[t.GetID()] {
		return
	}
	visited[t.GetID()] = true

	for _, node := range t.GetInputs() {
		_, ok := acc[node]
		if !ok {
			acc[node] = make([]INode, 0)
		}
		acc[node] = append(acc[node], t)

		getDepending(node, acc, visited)
	}
}
