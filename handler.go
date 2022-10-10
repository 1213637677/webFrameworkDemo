package webframework

type Routable interface {
	// Route 定义路由并执行 handleFunc 方法
	Route(method string, pattern string, handleFunc handlerFunc) error
}

type Handler interface {
	ServeHTTP(c *Context)
	Routable
}
