package webframework

import (
	"context"
	"fmt"
	"net/http"
)

type Server interface {
	Routable
	// Start 启动服务
	Start(address string)
	Shutdown(c context.Context) error
}

type SdkHttpServer struct {
	Name    string
	handler Handler
	root    Filter
}

// Route 定义路由并执行 handleFunc 方法
func (s *SdkHttpServer) Route(method string, pattern string, handleFunc handlerFunc) error {
	return s.handler.Route(method, pattern, handleFunc)
}

// Start 启动服务
func (s *SdkHttpServer) Start(address string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r)
		s.root(c)
	})
	http.ListenAndServe(address, nil)
}

func (s *SdkHttpServer) Shutdown(c context.Context) error {
	fmt.Println("server begin shutdown")
	// server shutdown 逻辑
	fmt.Println("server shutdown finish")
	return nil
}

func NewHttpServer(name string, builders ...FilterBuilder) Server {
	handler := NewHandlerBasedOnTree()
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
