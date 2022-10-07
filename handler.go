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
	http.Handler
	Routable
}

type HandlerBasedOnMap struct {
	handlers map[string]func(*Context)
}

func (h *HandlerBasedOnMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := h.key(r.Method, r.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		ctx := NewContext(w, r)
		handler(ctx)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found pattern"))
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
