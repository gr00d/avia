package application

import (
	"aviasales/internal/services"
	"aviasales/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type router struct {
	ginRouter *gin.Engine
}

func NewRouter(
	rootCtx context.Context,
	services services.IServiceFactory,
) *router {
	gin.SetMode(gin.ReleaseMode)

	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())
	pprof.Register(ginRouter)

	for i := range routes {
		currentRoute := routes[i]
		handler := ginRouter.GET
		if currentRoute.method == http.MethodPost {
			handler = ginRouter.POST
		}
		handler(currentRoute.path, func(c *gin.Context) {
			ctx := logger.With(rootCtx, "path", c.Request.URL.Path, "query", c.Request.URL.Query())
			timeOnStart := time.Now()

			currentRoute.handler.Process(c, services)

			logger.Info(ctx, "request processed", "duration", time.Since(timeOnStart).Seconds())
		})
	}

	ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger.json")))
	ginRouter.GET("swagger.json", func(c *gin.Context) {
		c.File("swagger.json")
	})
	ginRouter.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "I'm ok",
		})
	})

	return &router{
		ginRouter: ginRouter,
	}
}

func (r *router) GetRouterHandler() http.Handler {
	return r.ginRouter
}
