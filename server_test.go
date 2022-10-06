package webframework_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"webframework"
)

func TestSdkHttpServerRequestBody(t *testing.T) {
	svr := webframework.NewHttpServer("test")

	getBody := func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "read body failed, error: %v", err)
			return
		}
		fmt.Fprintf(w, "read body: %s", string(body))
	}

	svr.Route("/body", getBody)
	svr.Start("127.0.0.1:10228")
}

func TestSdkHttpServerRequestQuery(t *testing.T) {
	svr := webframework.NewHttpServer("test")
	getQuery := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		fmt.Fprintf(w, "query: %v", query)
	}

	svr.Route("/query", getQuery)
	svr.Start("127.0.0.1:10228")
}

func TestSdkHttpServerRequestHeader(t *testing.T) {
	svr := webframework.NewHttpServer("test")

	getHeader := func(w http.ResponseWriter, r *http.Request) {
		header := r.Header
		fmt.Fprintf(w, "header: %v", header)
	}

	svr.Route("/header", getHeader)
	svr.Start("127.0.0.1:10228")
}
