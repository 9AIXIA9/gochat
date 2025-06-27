package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gochat/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const defaultConfigPath = "./etc/config.yaml"

func main() {
	path := flag.String("config", defaultConfigPath, "config path")
	flag.Parse()

	// 使用Wire初始化依赖
	deps := InitializeDependencies(*path)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", deps.Config.Host, deps.Config.Port),
		Handler: api.Setup(deps),
	}

	// 启动服务器
	go func() {
		log.Printf("start server at : %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("start server failed: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("closing server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//todo 此处不会处理websocket连接 需自行处理
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("close server failed:", err)
	}
	log.Println("server has been closed")
}
