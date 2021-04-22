// Aviasales API
//
// Documentation for JAviasales API
//
// Version: 1.0.0
// swagger:meta
package main

import (
	"aviasales/internal/application"
	"aviasales/internal/services"
	"aviasales/internal/services/storage"
	"aviasales/pkg/logger"
	"aviasales/pkg/logger/zaplogger"
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const (
	fixturesDirectory = "./fixtures"
	parsersLimit      = 5
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	zapl, err := zaplogger.NewProduction()
	if err != nil {
		logger.FatalE(ctx, "unable to create logger", err)
	}
	defer zapl.Sync()
	logger.SetGlobalLogger(zapl)

	var logLevel logger.Level
	err = logLevel.Set("debug")
	if err != nil {
		logger.FatalE(ctx, "unable to parse loglevel", err)
	}
	logger.SetLevel(logLevel)

	serviceFactory := services.NewServiceFactory(ctx)
	parserWorkerPool := application.SpawnWorkers(parsersLimit)

	err = loadData(ctx, serviceFactory.Storage(), parserWorkerPool)
	if err != nil {
		logger.FatalE(ctx, "unable to load xml files", err)
	}

	server := application.NewServer(ctx, serviceFactory)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(quit)
		<-quit
		cancel()
		parserWorkerPool.Flush()
	}()

	server.Run()
}

func loadData(ctx context.Context, storage storage.IStorage, pool *application.WorkerPool) error {
	files, err := ioutil.ReadDir(fixturesDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := filepath.Join(fixturesDirectory, file.Name())
		pool.Put(&application.WorkerParserQueue{
			Ctx:      ctx,
			FileName: fileName,
			Storage:  storage,
		})
	}

	return nil
}
