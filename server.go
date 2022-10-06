package webframework

import "net/http"

type Server interface {
	// Route 定义路由并执行 handleFunc 方法
	Route(pattern string, handleFunc http.HandlerFunc)
	// Start 启动服务
	Start(address string)
}

type SdkHttpServer struct {
	Name string
}

// Route 定义路由并执行 handleFunc 方法
func (s *SdkHttpServer) Route(pattern string, handleFunc http.HandlerFunc) {
	http.Handle(pattern, handleFunc)
}

// Start 启动服务
func (s *SdkHttpServer) Start(address string) {
	http.ListenAndServe(address, nil)
}

func NewHttpServer(name string) Server {
	return &SdkHttpServer{
		Name: name,
	}
}
