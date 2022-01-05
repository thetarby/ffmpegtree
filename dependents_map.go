package ffmpegtree

// DependentsMap is a structure which maps nodes to all of its dependent nodes.
type DependentsMap struct{
	m map[INode][]INode
    keys []INode
}

func (m *DependentsMap) Set(k INode, val []INode){
	_, ok := m.m[k]
	m.m[k] = val
	if !ok{
		m.keys = append(m.keys, k)
	}
}

func (m *DependentsMap) GetAndCheck(k INode) ([]INode, bool){
	val, ok := m.m[k]
	return val, ok
}

func (m *DependentsMap) Get(k INode) []INode{
	return m.m[k]
}

func (m *DependentsMap) Keys() []INode{
	return m.keys
}

func NewDependentsMap() *DependentsMap{
	return &DependentsMap{
		m:    make(map[INode][]INode),
		keys: make([]INode, 0),
	}
}
