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
	defaultConfigPath = "./config/config.yaml"
)

func main() {
	path := flag.String("config", defaultConfigPath, "config file path")
	env := flag.String("env", "development", "environment (development, production)")
	flag.Parse()

	dependencies, err := initializeDependencies(*path, *env)
	if err != nil {
		log.Fatalf("initialize dependencies failed,err:%v", err)
	}

	// 创建HTTP服务器（用于 HTTPS）
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", dependencies.Config.Host, dependencies.Config.Port),
		Handler: setupRoutes(dependencies.Config, dependencies),
	}

	// 启动 HTTPS 服务器
	go func() {
		log.Printf("starting https server on %s", srv.Addr)
		// ListenAndServeTLS 会阻塞直到服务器返回错误
		if err := srv.ListenAndServeTLS(dependencies.Config.Cert.HTTPSCertFile, dependencies.Config.Cert.HTTPSKeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("start https server failed: %s\n", err)
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
	log.Println("server has been closed")
}
