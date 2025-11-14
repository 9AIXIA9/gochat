package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	defaultConfigFilePath = "./config/config.yaml"
	defaultENVFilePath    = "./.env.development"
)

func main() {
	path := flag.String("config", defaultConfigFilePath, "config file path")
	env := flag.String("env", defaultENVFilePath, "env file path")
	needMigrate := flag.Bool("need_migrate", false, "need migrate database or not")
	flag.Parse()

	dependencies, err := initializeDependencies(ConfigPath(*path), EnvPath(*env), *needMigrate)
	if err != nil {
		log.Fatalf("initialize Dependencies failed,err:%v", err)
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", dependencies.config.Host, dependencies.config.Port),
		Handler: dependencies.HttpRouter,
	}

	// 启动 HTTP 服务器
	go func() {
		log.Printf("starting http server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("start http server failed: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("closing server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("close server failed:", err)
	}
	if dependencies.closeAll != nil {
		dependencies.closeAll()
	}
	log.Println("server has been closed")
}
