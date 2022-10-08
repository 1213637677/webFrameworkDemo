package webframework_test

import (
	"fmt"
	"net/http"
	"testing"
	"webframework"
)

func TestHandlerBaseOnTree(t *testing.T) {
	svr := webframework.NewHttpServer("test")
	user := func(ctx *webframework.Context) {
		fmt.Printf("route /user/*\n")
	}
	user1 := func(ctx *webframework.Context) {
		fmt.Printf("route /user/1")
	}
	err := svr.Route(http.MethodGet, "/user/*", user)
	if err != nil {
		panic(err)
	}
	err = svr.Route(http.MethodGet, "/user/1", user1)
	if err != nil {
		panic(err)
	}
	svr.Start(":10228")
}
