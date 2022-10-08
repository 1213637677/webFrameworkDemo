package webframework

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Routable interface {
	// Route 定义路由并执行 handleFunc 方法
	Route(method string, pattern string, handleFunc handlerFunc) error
}

type Handler interface {
	ServeHTTP(c *Context)
	Routable
}

type HandlerBasedOnMap struct {
	handlers map[string]handlerFunc
}

func (h *HandlerBasedOnMap) ServeHTTP(c *Context) {
	key := h.key(c.R.Method, c.R.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		handler(c)
		return
	}
	c.W.WriteHeader(http.StatusNotFound)
	c.W.Write([]byte("not found pattern"))
}

func (h *HandlerBasedOnMap) key(method string, pattern string) string {
	return fmt.Sprintf("%s#%s", method, pattern)
}

// Route 定义路由并执行 handleFunc 方法
func (h *HandlerBasedOnMap) Route(method string, pattern string, handleFunc handlerFunc) error {
	key := h.key(method, pattern)
	h.handlers[key] = handleFunc
	return nil
}

func NewHandlerBasedOnMap() Handler {
	return &HandlerBasedOnMap{
		handlers: make(map[string]handlerFunc),
	}
}

type HandlerBasedOnTree struct {
	root *node
}

type handlerFunc func(c *Context)

type node struct {
	path     string
	children []*node

	handlerFunc handlerFunc
}

func (h *HandlerBasedOnTree) ServeHTTP(c *Context) {
	url := strings.Trim(c.R.URL.Path, "/")
	paths := strings.Split(url, "/")
	cur := h.root
	for _, path := range paths {
		if child, ok := cur.findMatchChild(path); ok {
			cur = child
			if child.path == "*" {
				break
			}
		} else {
			c.W.WriteHeader(http.StatusNotFound)
			c.W.Write([]byte("not found"))
			return
		}
	}
	if cur.handlerFunc == nil {
		c.W.WriteHeader(http.StatusNotFound)
		c.W.Write([]byte("not found"))
	} else {
		cur.handlerFunc(c)
	}
}

// Route 定义路由并执行 handleFunc 方法
func (h *HandlerBasedOnTree) Route(method string, pattern string, handleFunc handlerFunc) error {
	err := h.validatePattern(pattern)
	if err != nil {
		return err
	}
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")
	cur := h.root
	for idx, path := range paths {
		matchChild, ok := cur.findMatchChild(path)
		if ok && matchChild.path != "*" {
			cur = matchChild
		} else {
			cur.createSubTree(paths[idx:], handleFunc)
			break
		}
	}
	return nil
}

var ErrorInvalidRouterPattern = errors.New("invalid router pattern")

func (h *HandlerBasedOnTree) validatePattern(pattern string) error {
	pos := strings.Index(pattern, "*")
	if pos > 0 {
		if pos != len(pattern)-1 {
			return ErrorInvalidRouterPattern
		}
		if pattern[pos-1] != '/' {
			return ErrorInvalidRouterPattern
		}
	}
	return nil
}

func (n *node) findMatchChild(path string) (*node, bool) {
	var wildcardNode *node
	for _, child := range n.children {
		if child.path == path && child.path != "*" {
			return child, true
		}
		if child.path == "*" {
			wildcardNode = child
		}
	}
	return wildcardNode, wildcardNode != nil
}

func (n *node) createSubTree(paths []string, handlerFunc handlerFunc) {
	cur := n
	for _, path := range paths {
		newNode := newNode(path)
		cur.children = append(cur.children, newNode)
		cur = newNode
	}
	cur.handlerFunc = handlerFunc
}

func newNode(path string) *node {
	return &node{
		path:     path,
		children: make([]*node, 0),
	}
}

func NewHandlerBasedOnTree() Handler {
	return &HandlerBasedOnTree{
		root: newNode(""),
	}
}
