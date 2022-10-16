package webframework

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Hook func(c context.Context) error

var ErrorHookTimeout = errors.New("the hook timeout")

func BuildCloseServerHook(services ...Server) Hook {
	return func(c context.Context) error {
		wg := sync.WaitGroup{}
		done := make(chan struct{})
		wg.Add(len(services))
		for _, service := range services {
			go func(s Server) {
				err := s.Shutdown(c)
				if err != nil {
					fmt.Printf("shutdown server failed, error: %v", err)
				}
				wg.Done()
			}(service)
		}
		go func() {
			wg.Wait()
			done <- struct{}{}
		}()
		select {
		case <-c.Done():
			fmt.Println("close server timeout")
			return ErrorHookTimeout
		case <-done:
			fmt.Println("close all servers")
		}
		return nil
	}
}
