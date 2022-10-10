package webframework

import (
	"errors"
	"net/http"
	"strings"
)

type HandlerBasedOnTree struct {
	root Node
}

type handlerFunc func(c *Context)

func (h *HandlerBasedOnTree) ServeHTTP(c *Context) {
	url := strings.Trim(c.R.URL.Path, "/")
	paths := strings.Split(url, "/")
	cur := h.root
	for _, path := range paths {
		if child := cur.FindMatchChild(path); child != nil {
			cur = child
			if child.GetNodeType() == nodeTypeAny {
				break
			}
		} else {
			c.W.WriteHeader(http.StatusNotFound)
			c.W.Write([]byte("not found"))
			return
		}
	}
	if cur.GetHandlerFunc() == nil {
		c.W.WriteHeader(http.StatusNotFound)
		c.W.Write([]byte("not found"))
	} else {
		cur.GetHandlerFunc()(c)
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
		matchChild := cur.FindMatchChild(path)
		if matchChild != nil && matchChild.GetNodeType() != nodeTypeAny {
			cur = matchChild
		} else {
			createSubTree(cur, paths[idx:], handleFunc)
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

func NewHandlerBasedOnTree() Handler {
	return &HandlerBasedOnTree{
		root: newRootNode(nil),
	}
}
