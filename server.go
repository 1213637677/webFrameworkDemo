package webframework

import "net/http"

type Server interface {
	// Route 定义路由并执行 handleFunc 方法
	Route(method string, pattern string, handleFunc func(ctx *Context))
	// Start 启动服务
	Start(address string)
}

type SdkHttpServer struct {
	Name    string
	handler *HandlerBasedOnMap
}

// Route 定义路由并执行 handleFunc 方法
func (s *SdkHttpServer) Route(method string, pattern string, handleFunc func(ctx *Context)) {
	key := s.handler.key(method, pattern)
	s.handler.handlers[key] = handleFunc
}

// Start 启动服务
func (s *SdkHttpServer) Start(address string) {
	http.ListenAndServe(address, s.handler)
}

func NewHttpServer(name string) Server {
	return &SdkHttpServer{
		Name: name,
		handler: &HandlerBasedOnMap{
			handlers: make(map[string]func(ctx *Context)),
		},
	}
}
