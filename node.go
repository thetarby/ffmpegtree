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

func GetDependents(t... INode) *DependentsMap {
	res := NewDependentsMap()
	visited := make(map[string]bool)
	for _, n := range t{
		getDependents(n, res, visited)
	}

	return res
}

func getDependents(t INode, acc *DependentsMap, visited map[string]bool) {
	if visited[t.GetID()] {
		return
	}
	visited[t.GetID()] = true

	for _, node := range t.GetInputs() {
		val, ok := acc.GetAndCheck(node)
		if !ok {
			acc.Set(node, make([]INode, 0))
			val, _ = acc.GetAndCheck(node)
		}
		val = append(val, t)
		acc.Set(node, val)

		getDependents(node, acc, visited)
	}
}
