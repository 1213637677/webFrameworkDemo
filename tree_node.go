package webframework

import "sort"

const (
	nodeTypeRoot = iota
	nodeTypeAny
	nodeTypeParam
	nodeTypeReg
	nodeTypeStatic
)

type Node interface {
	// Match 当前节点是否匹配 path
	Match(path string) bool
	// GetNodeType 获取节点类型
	GetNodeType() int
	// FindMatchChild 查找子节点
	FindMatchChild(path string) Node
	// GetHandlerFunc 获取当前节点的 handlerFunc
	GetHandlerFunc() handlerFunc
	// SetHandlerFunc 添加 handlerFunc
	SetHandlerFunc(handlerFunc handlerFunc)
	// AddChildNode 添加子节点
	AddChildNode(child Node)
}

type baseNode struct {
	nodeType    int
	handlerFunc handlerFunc
	children    []Node
}

func (n *baseNode) GetNodeType() int {
	return n.nodeType
}

func (n *baseNode) FindMatchChild(path string) Node {
	candidates := make([]Node, 0, 2)
	for _, child := range n.children {
		if child.Match(path) {
			candidates = append(candidates, child)
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].GetNodeType() > candidates[j].GetNodeType()
	})
	return candidates[0]
}

func (n *baseNode) GetHandlerFunc() handlerFunc {
	return n.handlerFunc
}

func (n *baseNode) SetHandlerFunc(handlerFunc handlerFunc) {
	n.handlerFunc = handlerFunc
}

func (n *baseNode) AddChildNode(child Node) {
	n.children = append(n.children, child)
}

func createSubTree(node Node, paths []string, handlerFunc handlerFunc) {
	cur := node
	for _, path := range paths {
		if path != "*" {
			child := newStaticNode(path, nil)
			cur.AddChildNode(child)
		}
	}
	cur.SetHandlerFunc(handlerFunc)
}

func newBaseNode(typ int, handlerFunc handlerFunc) baseNode {
	return baseNode{
		nodeType:    typ,
		handlerFunc: handlerFunc,
		children:    make([]Node, 0, 2),
	}
}

type staticNode struct {
	baseNode
	path string
}

func (n *staticNode) Match(path string) bool {
	return n.path == path
}

func newStaticNode(path string, handlerFunc handlerFunc) Node {
	return &staticNode{
		baseNode: newBaseNode(nodeTypeStatic, handlerFunc),
		path:     path,
	}
}

type rootNode struct {
	baseNode
}

func (n *rootNode) Match(path string) bool {
	return true
}

func newRootNode(handlerFunc handlerFunc) Node {
	return &rootNode{
		baseNode: newBaseNode(nodeTypeRoot, handlerFunc),
	}
}
