package main

import (
	"github.com/gin-gonic/gin"
	authorizationRouter "gochat/internal/authorization/ports/http/routers"
	"gochat/internal/shared/config"
	"gochat/internal/shared/http/middlewares"
)

func setupRoutes(appConfig *config.App, dependencies *Dependencies) *gin.Engine {
	engine := gin.New()

	engine.Use(middlewares.Logger(), middlewares.Recover(), middlewares.CORS(appConfig.CORS))

	baseGroup := engine.Group("/api/v1")

	authorizationRouter.Setup(baseGroup.Group("/authorization"), dependencies.AuthorizationDependencies)

	return engine
}
