package storage

import (
	"aviasales/pkg/entities"
	"context"
)

type IStorage interface {
	AddItinerary(itinerary entities.Itinerary)
	GetItineraries(start, destination string) ([]*entities.Itinerary, error)
	GetCheapest(start, destination string) (*entities.Itinerary, error)
	GetMostExpensive(start, destination string) (*entities.Itinerary, error)
	GetLongest(start, destination string) (*entities.Itinerary, error)
	GetShortest(start, destination string) (*entities.Itinerary, error)
	GetOptimal(start, destination string) (*entities.Itinerary, error)
	GetByUUID(UUID string) (*entities.Itinerary, error)
}

func New(ctx context.Context) *service {
	return NewMemoryStorage(ctx)
}
