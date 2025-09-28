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
const defaultEnvPath = "development"

func main() {
	path := flag.String("config", defaultConfigPath, "config path")
	env := flag.String("env", defaultEnvPath, "environment (development, production, etc)")
	flag.Parse()

	// 如果命令行指定了环境，设置到环境变量中
	if err := os.Setenv("GOCHAT_ENV", *env); err != nil {
		log.Fatalf("set env failed,err:%v", err)
	}
	log.Printf("using environment: %s", *env)

	// 使用Wire初始化依赖
	deps := InitializeDependencies(*path)

	// 读取证书路径（可通过环境变量覆盖）
	certFile := os.Getenv("TLS_CERT")
	if certFile == "" {
		certFile = "cert.pem"
	}
	keyFile := os.Getenv("TLS_KEY")
	if keyFile == "" {
		keyFile = "key.pem"
	}

	// 确认证书文件存在
	if _, err := os.Stat(certFile); err != nil {
		log.Fatalf("tls cert file not found: %s (generate with mkcert or set TLS_CERT). err: %v", certFile, err)
	}
	if _, err := os.Stat(keyFile); err != nil {
		log.Fatalf("tls key file not found: %s (generate with mkcert or set TLS_KEY). err: %v", keyFile, err)
	}

	// 创建HTTP服务器（用于 HTTPS）
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", deps.Config.Host, deps.Config.Port),
		Handler: api.Setup(deps),
	}

	// 启动 HTTPS 服务器
	go func() {
		log.Printf("starting https server on %s (cert=%s key=%s)", srv.Addr, certFile, keyFile)
		// ListenAndServeTLS 会阻塞直到服务器返回错误
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
