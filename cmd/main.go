// @title           GoChat Backend API
// @version         1.0
// @description     GoChat 聊天/好友/房间管理后端接口
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.email   support@example.com
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
// @host            localhost:8888
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
package main

import (
	"context"
	"flag"
	"gochat/cmd/di"
	"gochat/internal/infrastructure/godotenv"
	otelInfra "gochat/internal/infrastructure/otel"
	"gochat/internal/infrastructure/viper"
	zaputils "gochat/internal/infrastructure/zap"
	"log"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
)

const (
	defaultConfigFilePath = "./config/config.yaml"
	defaultENVFilePath    = "./.env.development"
)

func main() {
	path := flag.String("config", defaultConfigFilePath, "config file path")
	env := flag.String("env", defaultENVFilePath, "env file path")
	flag.Parse()

	if err := godotenv.LoadEnvFile(*env); err != nil {
		log.Fatalf("load env failed,err:%v", err)
	}

	conf, err := viper.LoadConfigFile(*path)
	if err != nil {
		log.Fatalf("load config file failed,err:%v", err)
	}

	closeOtel, teleErr := otelInfra.Initialize(context.Background(), conf.Name, conf.Telemetry)
	if teleErr != nil {
		zap.L().Warn("initialize telemetry failed", zap.Error(teleErr))
	}

	if err := zaputils.Initialize(conf.Logger); err != nil {
		log.Fatalf("initialize zap failed:%v", err)
	}

	dependencies, err := di.Initialize(conf)
	if err != nil {
		log.Fatalf("initialize dependencies failed,err:%v", err)
	}

	zap.L().Info("server started", zap.String("app", conf.Name))

	//启动各个组件
	dependencies.EmailNotifier.Start()
	dependencies.KafkaEventPublisher.Start()
	for _, c := range dependencies.KafkaConsumers {
		if c == nil { // safety
			continue
		}
		if err := c.Start(); err != nil {
			zap.L().Fatal("start kafka event subscriber failed", zap.Error(err))
		}
	}
	dependencies.BinlogReader.Start()
	dependencies.HttpServer.Start()

	//优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	zap.L().Info("server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if err := dependencies.HttpServer.Close(ctx); err != nil {
		zap.L().Error("shutdown http server failed", zap.Error(err))
	}

	dependencies.EmailNotifier.Close()
	dependencies.KafkaEventPublisher.Close()
	for i := len(dependencies.KafkaConsumers) - 1; i >= 0; i-- {
		if dependencies.KafkaConsumers[i] != nil {
			dependencies.KafkaConsumers[i].Close()
		}
	}
	dependencies.BinlogReader.Close()

	if closeOtel != nil {
		if err := closeOtel(ctx); err != nil {
			zap.L().Error("shutdown telemetry failed", zap.Error(err))
		}
	}

	zap.L().Info("server exited properly")
}
