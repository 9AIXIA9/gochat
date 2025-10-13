package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gochat/api"
	"gochat/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const defaultConfigPath = "./etc/config.yaml"
const defaultEnvPath = "development"

func main() {
	path := flag.String("config", defaultConfigPath, "config path")
	env := flag.String("env", defaultEnvPath, "environment (development, production, etc)")
	flag.Parse()

	// 如果命令行指定了环境，设置到环境变量中
	if err := os.Setenv("GOCHAT_ENV", *env); err != nil {
		log.Fatalf("set env failed,err:%v", err)
	}
	if err := config.LoadEnvFile(); err != nil {
		log.Fatalf("load env failed,err:%v", err)
	}
	log.Printf("using environment: %s", *env)

	// 使用Wire初始化依赖
	deps, err := NewDependencies(*path)
	if err != nil {
		log.Fatalf("init dependencies failed,err:%v", err)
	}

	// 确认证书文件存在
	if _, err := os.Stat(deps.Config.Cert.HTTPSCertFile); err != nil {
		log.Fatalf("tls cert file not found: %s (generate with mkcert or set TLS_CERT). err: %v", deps.Config.Cert.HTTPSCertFile, err)
	}
	if _, err := os.Stat(deps.Config.Cert.HTTPSKeyFile); err != nil {
		log.Fatalf("tls key file not found: %s (generate with mkcert or set TLS_KEY). err: %v", deps.Config.Cert.HTTPSKeyFile, err)
	}

	// 创建HTTP服务器（用于 HTTPS）
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", deps.Config.Host, deps.Config.Port),
		Handler: api.Setup(deps),
	}

	// 启动 HTTPS 服务器
	go func() {
		log.Printf("starting https server on %s", srv.Addr)
		// ListenAndServeTLS 会阻塞直到服务器返回错误
		if err := srv.ListenAndServeTLS(deps.Config.Cert.HTTPSCertFile, deps.Config.Cert.HTTPSKeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
