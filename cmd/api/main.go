// @title           GoChat Backend API
// @version         1.0
// @description     GoChat 聊天/好友/房间管理后端接口
// @termsOfService  https://swagger.io/terms/
// @contact.name    XIA
// @contact.email   906094554@qq.com
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
// @BasePath        /api/v1
// @schemes         http
// @accept json
// @produce json
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 认证格式：Bearer {access_token}（注意 Bearer 后加英文空格）
//
//go:generate swag init -g ./cmd/api/main.go -o ../../docs --parseDependency --parseInternal --dir ../..
package main

import (
	"context"
	"flag"
	"fmt"
	"gochat/cmd/api/di"
	"gochat/internal/infrastructure/godotenv"
	"gochat/internal/infrastructure/viper"
	zaputils "gochat/internal/infrastructure/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

const (
	defaultConfigFilePath = "./config/config.yaml"
	shutdownTimeout       = 30 * time.Second
)

func main() {
	baseConfigPath := flag.String("config", defaultConfigFilePath, "config file path")
	overrideConfigPath := flag.String("override_config", "", "override config file path")

	env := flag.String("env", "", "env file path")
	flag.Parse()

	if env != nil && *env != "" {
		if err := godotenv.LoadEnvFile(*env); err != nil {
			log.Fatalf("load env failed,err:%v", err)
		}
	}

	conf, err := viper.LoadConfigFile(*baseConfigPath, *overrideConfigPath)
	if err != nil {
		log.Fatalf("load config file failed,err:%v", err)
	}

	if err := zaputils.Initialize(conf.Logger); err != nil {
		log.Fatalf("initialize zap failed:%v", err)
	}

	dependencies, err := di.Initialize(conf)
	if err != nil {
		log.Fatalf("initialize dependencies failed,err:%v", err)
	}

	if err := startComponents(dependencies); err != nil {
		zap.L().Error("start components failed", zap.Error(err))
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		shutdownComponents(shutdownCtx, dependencies)
		log.Fatalf("server startup failed, rolled back: %v", err)
	}

	zap.L().Info("server started", zap.String("app", conf.Name))

	// Use context-based signal handling for one-shot graceful shutdown.
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-sigCtx.Done()

	zap.L().Info("server is shutting down...", zap.Duration("timeout", shutdownTimeout))

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	shutdownComponents(ctx, dependencies)
	zap.L().Info("server exited properly")
}

func startComponents(dependencies *di.Dependencies) error {
	dependencies.EmailNotifier.Start()
	dependencies.KafkaEventPublisher.Start()

	startedConsumers := 0
	for _, c := range dependencies.KafkaConsumers {
		if c == nil {
			continue
		}
		if err := c.Start(); err != nil {
			for i := startedConsumers - 1; i >= 0; i-- {
				if dependencies.KafkaConsumers[i] != nil {
					dependencies.KafkaConsumers[i].Close()
				}
			}
			dependencies.KafkaEventPublisher.Close()
			dependencies.EmailNotifier.Close()
			return fmt.Errorf("start kafka event subscriber failed: %w", err)
		}
		startedConsumers++
	}

	dependencies.BinlogReader.Start()
	dependencies.HttpServer.Start()
	return nil
}

func shutdownComponents(ctx context.Context, dependencies *di.Dependencies) {
	if err := dependencies.HttpServer.Close(ctx); err != nil {
		zap.L().Error("shutdown http server failed", zap.Error(err))
	}

	dependencies.BinlogReader.Close()

	for i := len(dependencies.KafkaConsumers) - 1; i >= 0; i-- {
		if dependencies.KafkaConsumers[i] != nil {
			dependencies.KafkaConsumers[i].Close()
		}
	}

	dependencies.KafkaEventPublisher.Close()
	dependencies.EmailNotifier.Close()

	if dependencies.OTELShutdown != nil {
		if err := dependencies.OTELShutdown(ctx); err != nil {
			zap.L().Error("shutdown otel failed", zap.Error(err))
		}
	}
}
