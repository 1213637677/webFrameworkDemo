package webframework

import (
	"fmt"
	"net/http"
)

type Routable interface {
	// Route 定义路由并执行 handleFunc 方法
	Route(method string, pattern string, handleFunc func(ctx *Context))
}

type Handler interface {
	ServeHTTP(c *Context)
	Routable
}

type HandlerBasedOnMap struct {
	handlers map[string]func(*Context)
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
func (h *HandlerBasedOnMap) Route(method string, pattern string, handleFunc func(ctx *Context)) {
	key := h.key(method, pattern)
	h.handlers[key] = handleFunc
}

func NewHandlerBasedOnMap() Handler {
	return &HandlerBasedOnMap{
		handlers: make(map[string]func(*Context)),
	}
}
