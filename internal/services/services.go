package services

import (
	"aviasales/internal/services/storage"
	"context"
	"sync"
)

type factory struct {
	ctx      context.Context
	safeInit servicesInitLocks
	storage  storage.IStorage
}

type servicesInitLocks struct {
	storage sync.Once
}

type IServiceFactory interface {
	Storage() storage.IStorage
}

func NewServiceFactory(
	ctx context.Context,
) IServiceFactory {
	return &factory{
		ctx: ctx,
	}
}

func (f *factory) Storage() storage.IStorage {
	f.safeInit.storage.Do(func() {
		f.storage = storage.New(f.ctx)
	})
	return f.storage
}
