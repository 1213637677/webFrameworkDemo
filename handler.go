package webframework

import (
	"fmt"
	"net/http"
)

type HandlerBasedOnMap struct {
	handlers map[string]func(ctx *Context)
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
