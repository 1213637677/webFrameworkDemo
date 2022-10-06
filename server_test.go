package webframework_test

import (
	"fmt"
	"io"
	"testing"
	"webframework"
)

func TestSdkHttpServerRequestBody(t *testing.T) {
	svr := webframework.NewHttpServer("test")

	getBody := func(ctx *webframework.Context) {
		body, err := io.ReadAll(ctx.R.Body)
		if err != nil {
			fmt.Fprintf(ctx.W, "read body failed, error: %v", err)
			return
		}
		fmt.Fprintf(ctx.W, "read body: %s", string(body))
	}

	svr.Route("/body", getBody)
	svr.Start("127.0.0.1:10228")
}

func TestSdkHttpServerRequestQuery(t *testing.T) {
	svr := webframework.NewHttpServer("test")
	getQuery := func(ctx *webframework.Context) {
		query := ctx.R.URL.Query()
		fmt.Fprintf(ctx.W, "query: %v", query)
	}

	svr.Route("/query", getQuery)
	svr.Start("127.0.0.1:10228")
}

func TestSdkHttpServerRequestHeader(t *testing.T) {
	svr := webframework.NewHttpServer("test")

	getHeader := func(ctx *webframework.Context) {
		header := ctx.R.Header
		fmt.Fprintf(ctx.W, "header: %v", header)
	}

	svr.Route("/header", getHeader)
	svr.Start("127.0.0.1:10228")
}
