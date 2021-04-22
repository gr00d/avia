package storage

import (
	"aviasales/pkg/entities"
	"aviasales/pkg/logger"
	"context"
	"errors"
	"sync"

	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type service struct {
	ctx            context.Context
	Itineraries    map[entities.SourceCity]map[entities.DestinationCity]*entities.Itineraries
	ItinerariesMap map[entities.ItineraryUUID]entities.Itinerary
	isUpdating     sync.RWMutex
}

func NewMemoryStorage(ctx context.Context) *service {
	itineraries := make(map[entities.SourceCity]map[entities.DestinationCity]*entities.Itineraries)
	return &service{
		ctx:            ctx,
		Itineraries:    itineraries,
		ItinerariesMap: map[entities.ItineraryUUID]entities.Itinerary{},
	}
}

func (s *service) AddItinerary(itinerary entities.Itinerary) {
	s.isUpdating.Lock()
	defer s.isUpdating.Unlock()

	if len(itinerary.Onward) == 0 {
		logger.Debug(s.ctx, "unable to find source point")
		return
	}

	sourcePoint := entities.SourceCity(itinerary.Onward[0].Source)
	_, ok := s.Itineraries[sourcePoint]
	if !ok {
		s.Itineraries[sourcePoint] = make(map[entities.DestinationCity]*entities.Itineraries)
	}

	destinationPoint := entities.DestinationCity(itinerary.Onward[len(itinerary.Onward)-1].Destination)
	_, ok = s.Itineraries[sourcePoint][destinationPoint]
	if !ok {
		s.Itineraries[sourcePoint][destinationPoint] = &entities.Itineraries{
			Itineraries:   []*entities.Itinerary{},
			Shortest:      nil,
			Longest:       nil,
			Cheapest:      nil,
			MostExpensive: nil,
			Optimal:       nil,
		}
	}

	itineraryUUID := entities.ItineraryUUID(uuid.NewV4().String())
	itinerary.UUID = itineraryUUID
	s.ItinerariesMap[itineraryUUID] = itinerary

	itineraries := s.Itineraries[sourcePoint][destinationPoint]
	itineraries.Itineraries = append(itineraries.Itineraries, &itinerary)
	itineraries.Cheapest = getCheapest(itineraries.Cheapest, &itinerary, entities.ChargeTypeTotalAmount, entities.TypeSingleAdult)
	itineraries.MostExpensive = getMostExpensiveCheapest(itineraries.MostExpensive, &itinerary, entities.ChargeTypeTotalAmount, entities.TypeSingleAdult)
	itineraries.Shortest = getShortest(itineraries.Shortest, &itinerary)
	itineraries.Longest = getLongest(itineraries.Longest, &itinerary)
	itineraries.Optimal = getOptimal(itineraries.Optimal, &itinerary)
}

func (s *service) GetItineraries(source, destination string) ([]*entities.Itinerary, error) {
	itineraries, ok := s.Itineraries[entities.SourceCity(source)][entities.DestinationCity(destination)]
	if !ok {
		return []*entities.Itinerary{}, nil
	}

	return itineraries.Itineraries, nil
}

func (s *service) GetCheapest(source, destination string) (*entities.Itinerary, error) {
	itineraries, ok := s.Itineraries[entities.SourceCity(source)][entities.DestinationCity(destination)]
	if !ok {
		return &entities.Itinerary{}, nil
	}

	return itineraries.Cheapest, nil
}

func (s *service) GetMostExpensive(source, destination string) (*entities.Itinerary, error) {
	itineraries, ok := s.Itineraries[entities.SourceCity(source)][entities.DestinationCity(destination)]
	if !ok {
		return &entities.Itinerary{}, nil
	}

	return itineraries.MostExpensive, nil
}

func (s *service) GetLongest(source, destination string) (*entities.Itinerary, error) {
	itineraries, ok := s.Itineraries[entities.SourceCity(source)][entities.DestinationCity(destination)]
	if !ok {
		return &entities.Itinerary{}, nil
	}

	return itineraries.Longest, nil
}

func (s *service) GetShortest(source, destination string) (*entities.Itinerary, error) {
	itineraries, ok := s.Itineraries[entities.SourceCity(source)][entities.DestinationCity(destination)]
	if !ok {
		return &entities.Itinerary{}, nil
	}

	return itineraries.Shortest, nil
}

func (s *service) GetOptimal(source, destination string) (*entities.Itinerary, error) {
	itineraries, ok := s.Itineraries[entities.SourceCity(source)][entities.DestinationCity(destination)]
	if !ok {
		return &entities.Itinerary{}, nil
	}

	return itineraries.Optimal, nil
}

func (s *service) GetByUUID(uuid string) (*entities.Itinerary, error) {
	itinerary, ok := s.ItinerariesMap[entities.ItineraryUUID(uuid)]
	if !ok {
		return nil, errors.New("unable to find ticket")
	}

	return &itinerary, nil
}

func getCheapest(itinerary1, itinerary2 *entities.Itinerary, chargeType, rateType string) *entities.Itinerary {
	if itinerary1 == nil {
		return itinerary2
	}

	if itinerary1.GetPrice(chargeType, rateType).GreaterThan(itinerary2.GetPrice(chargeType, rateType)) {
		return itinerary2
	}

	return itinerary1
}

func getMostExpensiveCheapest(itinerary1, itinerary2 *entities.Itinerary, chargeType, rateType string) *entities.Itinerary {
	if itinerary1 == nil {
		return itinerary2
	}

	if itinerary1.GetPrice(chargeType, rateType).GreaterThan(itinerary2.GetPrice(chargeType, rateType)) {
		return itinerary1
	}

	return itinerary2
}

func getShortest(itinerary1, itinerary2 *entities.Itinerary) *entities.Itinerary {
	if itinerary1 == nil {
		return itinerary2
	}

	if itinerary1.GetDuration() < itinerary2.GetDuration() {
		return itinerary1
	}

	return itinerary2
}

func getLongest(itinerary1, itinerary2 *entities.Itinerary) *entities.Itinerary {
	if itinerary1 == nil {
		return itinerary2
	}

	if itinerary1.GetDuration() > itinerary2.GetDuration() {
		return itinerary1
	}

	return itinerary2
}

func getOptimal(itinerary1, itinerary2 *entities.Itinerary) *entities.Itinerary {
	if itinerary1 == nil {
		return itinerary2
	}

	score1, score2 := 0, 0
	var coefficient1 decimal.Decimal
	var coefficient2 decimal.Decimal
	if itinerary1.GetPrice(entities.ChargeTypeTotalAmount, entities.TypeSingleAdult).GreaterThan(decimal.NewFromInt(0)) {
		coefficient1 = decimal.NewFromInt(itinerary1.GetDurationWithoutTransfer()).Div(itinerary1.GetPrice(entities.ChargeTypeTotalAmount, entities.TypeSingleAdult))
	}
	if itinerary2.GetPrice(entities.ChargeTypeTotalAmount, entities.TypeSingleAdult).GreaterThan(decimal.NewFromInt(0)) {
		coefficient2 = decimal.NewFromInt(itinerary2.GetDurationWithoutTransfer()).Div(itinerary2.GetPrice(entities.ChargeTypeTotalAmount, entities.TypeSingleAdult))
	}

	if coefficient1.GreaterThan(coefficient2) {
		score1++
	} else {
		score2++
	}

	if itinerary1.GetPrice(entities.ChargeTypeTotalAmount, entities.TypeSingleAdult).
		LessThan(itinerary2.GetPrice(entities.ChargeTypeTotalAmount, entities.TypeSingleAdult)) {
		score1++
	} else {
		score2++
	}

	if itinerary1.GetDurationWithoutTransfer() < itinerary2.GetDurationWithoutTransfer() {
		score1++
	} else {
		score2++
	}

	if itinerary1.GetTransferDuration() < itinerary2.GetTransferDuration() {
		score1++
	} else {
		score2++
	}

	if score1 > score2 {
		return itinerary1
	}

	return itinerary2
}
