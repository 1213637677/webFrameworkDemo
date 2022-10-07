package webframework

import "net/http"

type Server interface {
	Routable
	// Start 启动服务
	Start(address string)
}

type SdkHttpServer struct {
	Name    string
	handler Handler
	root    Filter
}

// Route 定义路由并执行 handleFunc 方法
func (s *SdkHttpServer) Route(method string, pattern string, handleFunc func(ctx *Context)) {
	s.handler.Route(method, pattern, handleFunc)
}

// Start 启动服务
func (s *SdkHttpServer) Start(address string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r)
		s.root(c)
	})
	http.ListenAndServe(address, nil)
}

func NewHttpServer(name string, builders ...FilterBuilder) Server {
	handler := NewHandlerBasedOnMap()
	var root Filter = handler.ServeHTTP
	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}
	return &SdkHttpServer{
		Name:    name,
		handler: handler,
		root:    root,
	}
}
