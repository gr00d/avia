package application

import (
	"aviasales/internal/services"
	"aviasales/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	HTTPPort        = 8080
	gracefulTimeOut = 10 * time.Second
)

type server struct {
	ctx    context.Context
	router *router
}

func NewServer(
	ctx context.Context,
	services services.IServiceFactory,
) *server {
	return &server{
		ctx:    ctx,
		router: NewRouter(ctx, services),
	}
}

func (s *server) Run() {
	ctx := logger.With(s.ctx, "app", "run")
	addr := fmt.Sprintf("0.0.0.0:%d", HTTPPort)

	srv := &http.Server{
		Addr:    addr,
		Handler: s.router.GetRouterHandler(),
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, "error while ListenAndServe", err)
		}
	}()

	// graceful shutdown
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-s.ctx.Done()
		logger.Info(s.ctx, "Shutting down http server...")
		sCtx, cancel := context.WithTimeout(context.Background(), gracefulTimeOut)
		defer func() {
			cancel()
		}()
		if err := srv.Shutdown(sCtx); err != nil {
			logger.Error(ctx, "error while shutdown server", err)
		}
	}()

	wg.Wait()
	logger.Info(ctx, "Server is shutdown")
}
