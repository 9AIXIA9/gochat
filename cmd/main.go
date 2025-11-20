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

	dependencies, err := di.Initialize(conf)
	if err != nil {
		log.Fatalf("initialize dependencies failed,err:%v", err)
	}

	closeOtel, teleErr := otelInfra.Initialize(context.Background(), conf.Name, conf.Telemetry)
	if teleErr != nil {
		zap.L().Warn("initialize telemetry failed", zap.Error(teleErr))
	}

	if err := zaputils.Initialize(conf.Logger); err != nil {
		log.Fatalf("initialize zap failed:%v", err)
	}

	dependencies.EmailNotifier.Start()
	dependencies.KafkaEventPublisher.Start()
	if err := dependencies.KafkaEventSubscriber.Start(conf.Name); err != nil {
		zap.L().Fatal("start kafka event subscriber failed", zap.Error(err))
	}
	dependencies.OutboxConsumer.Start()
	dependencies.HttpServer.Start()

	//优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if err := dependencies.HttpServer.Close(ctx); err != nil {
		zap.L().Error("shutdown http server failed", zap.Error(err))
	}

	dependencies.EmailNotifier.Close()
	dependencies.KafkaEventPublisher.Close()
	dependencies.KafkaEventSubscriber.Close()
	dependencies.OutboxConsumer.Close()

	if closeOtel != nil {
		if err := closeOtel(ctx); err != nil {
			zap.L().Error("shutdown telemetry failed", zap.Error(err))
		}
	}

	zap.L().Info("server exited properly")
}
