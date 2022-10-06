package webframework_test

import (
	"fmt"
	"net/http"
	"testing"
	"webframework"
)

type user struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestContext(t *testing.T) {
	svr := webframework.NewHttpServer("test")
	userMap := make(map[string]*user)
	createUser := func(ctx *webframework.Context) {
		fmt.Println("start")
		newUser := &user{}
		err := ctx.ReadJson(newUser)
		if err != nil {
			fmt.Printf("read json failed, error: %v\n", err)
			ctx.WriteJson(http.StatusBadRequest, err)
			return
		}
		userMap[newUser.Name] = newUser
		err = ctx.WriteJson(http.StatusOK, newUser)
		if err != nil {
			fmt.Printf("write json failed, error: %v\n", err)
			ctx.WriteJson(http.StatusBadRequest, err)
		}
	}
	svr.Route("/user/create", createUser)
	svr.Start(":10228")
}
