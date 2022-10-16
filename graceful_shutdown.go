package webframework

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

type GracefulShutdown struct {
	// 正在处理中的请求数
	reqCnt int64
	// 关闭参数
	closing int32
	// 通知所有请求处理完毕
	zeroReqCnt chan struct{}
}

func NewGracefulShutdown() *GracefulShutdown {
	return &GracefulShutdown{
		zeroReqCnt: make(chan struct{}),
	}
}

// ShutdownFilterBuilder 关闭服务时，拒绝请求。当服务中请求处理完毕时，向 zeroReqCnt 中塞入数据
func (gs *GracefulShutdown) ShutdownFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		cl := atomic.LoadInt32(&gs.closing)
		if cl > 0 {
			c.W.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		atomic.AddInt64(&gs.reqCnt, 1)
		next(c)
		n := atomic.AddInt64(&gs.reqCnt, -1)
		cl = atomic.LoadInt32(&gs.closing)
		if cl > 0 && n == 0 {
			gs.zeroReqCnt <- struct{}{}
		}
	}
}

func (gs *GracefulShutdown) RejectNewRequestAndWaiting(ctx context.Context) error {
	atomic.AddInt32(&gs.closing, 1)
	if atomic.LoadInt64(&gs.reqCnt) == 0 {
		return nil
	}
	done := ctx.Done()
	select {
	case <-done:
		fmt.Println("超时")
	case <-gs.zeroReqCnt:
		fmt.Println("请求处理完毕")
	}
	return nil
}

func WaitForShutdown(hooks ...Hook) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, ShutdownSignals...)
	sig := <-signals
	fmt.Printf("get signal %s, application will shotdown\n", sig)
	time.AfterFunc(time.Minute*10, func() {
		fmt.Printf("Shutdown gracefully timeout, application will shutdown immediately. ")
		os.Exit(1)
	})
	for _, hook := range hooks {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		err := hook(ctx)
		if err != nil {
			fmt.Printf("failed to run hook, err: %v\n", err)
		}
		cancel()
	}
	os.Exit(1)
}
