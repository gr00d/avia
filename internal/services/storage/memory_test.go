package storage

import (
	"aviasales/pkg/entities"
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

const (
	source          = "DXB"
	destination     = "BKK"
	decimalEqualNum = 0
)

var itineraries = []entities.Itinerary{
	entities.Itinerary{
		UUID: "",
		Onward: []entities.Flight{
			entities.Flight{
				Carrier:            "",
				FlightNumber:       "",
				Source:             "DXB",
				Destination:        "DEL",
				DepartureTimeStamp: entities.FlightDate{Time: time.Date(2018, 10, 27, 00, 00, 00, 00, time.UTC)}, // 2018-10-27T0000
				ArrivalTimeStamp:   entities.FlightDate{Time: time.Date(2018, 10, 27, 04, 00, 00, 00, time.UTC)}, // 2018-10-27T0400
				Class:              "",
				NumberOfStops:      "",
				TicketType:         "",
			},
			entities.Flight{
				Carrier:            "",
				FlightNumber:       "",
				Source:             "DEL",
				Destination:        "BKK",
				DepartureTimeStamp: entities.FlightDate{Time: time.Date(2018, 10, 27, 13, 00, 00, 00, time.UTC)}, // 2018-10-27T1300
				ArrivalTimeStamp:   entities.FlightDate{Time: time.Date(2018, 10, 27, 19, 00, 00, 00, time.UTC)}, // 2018-10-27T1900
				Class:              "",
				NumberOfStops:      "",
				TicketType:         "",
			},
		},
		Return: nil,
		Pricing: &entities.Price{
			Currency: "SGD",
			ServiceCharges: []entities.Charge{
				{
					ChargeType: "TotalAmount",
					Type:       "BaseFare",
					Cost:       decimal.NewFromFloat(233.00),
				},
				{
					ChargeType: "TotalAmount",
					Type:       "AirlineTaxes",
					Cost:       decimal.NewFromFloat(152.40),
				},
				{
					ChargeType: "TotalAmount",
					Type:       "SingleAdult",
					Cost:       decimal.NewFromFloat(385.40),
				},
				{
					ChargeType: "TotalAmount",
					Type:       "SingleInfant",
					Cost:       decimal.NewFromFloat(140.40),
				},
			},
		},
	},
	entities.Itinerary{
		UUID: "",
		Onward: []entities.Flight{
			entities.Flight{
				Carrier:            "",
				FlightNumber:       "",
				Source:             "DXB",
				Destination:        "CAN",
				DepartureTimeStamp: entities.FlightDate{Time: time.Date(2018, 10, 27, 01, 00, 00, 00, time.UTC)}, // 2018-10-27T0100
				ArrivalTimeStamp:   entities.FlightDate{Time: time.Date(2018, 10, 27, 12, 00, 00, 00, time.UTC)}, // 2018-10-27T1200
				Class:              "",
				NumberOfStops:      "",
				TicketType:         "",
			},
			entities.Flight{
				Carrier:            "",
				FlightNumber:       "",
				Source:             "CAN",
				Destination:        "BKK",
				DepartureTimeStamp: entities.FlightDate{Time: time.Date(2018, 10, 27, 14, 00, 00, 00, time.UTC)}, // 2018-10-27T1400
				ArrivalTimeStamp:   entities.FlightDate{Time: time.Date(2018, 10, 27, 17, 00, 00, 00, time.UTC)}, // 2018-10-27T1700
				Class:              "",
				NumberOfStops:      "",
				TicketType:         "",
			},
		},
		Return: nil,
		Pricing: &entities.Price{
			Currency: "SGD",
			ServiceCharges: []entities.Charge{
				{
					ChargeType: "BaseFare",
					Type:       "SingleAdult",
					Cost:       decimal.NewFromFloat(167.00),
				},
				{
					ChargeType: "AirlineTaxes",
					Type:       "SingleAdult",
					Cost:       decimal.NewFromFloat(215.70),
				},
				{
					ChargeType: "TotalAmount",
					Type:       "SingleAdult",
					Cost:       decimal.NewFromFloat(382.70),
				},
			},
		},
	},
}

func TestService_GetCheapest(t *testing.T) {
	items := map[string]struct {
		expectedValue decimal.Decimal
	}{
		"it should be 382.70": {
			expectedValue: decimal.NewFromFloat(382.70),
		},
	}

	storage := NewMemoryStorage(context.Background())
	for i := range itineraries {
		storage.AddItinerary(itineraries[i])
	}

	for message, item := range items {
		itinerary, _ := storage.GetCheapest(source, destination)
		cmp := itinerary.GetPrice(entities.ChargeTypeTotalAmount, entities.TypeSingleAdult).Cmp(item.expectedValue)
		assert.Equal(t, decimalEqualNum, cmp, message)
	}
}

func TestService_GetMostExpensive(t *testing.T) {
	items := map[string]struct {
		expectedValue decimal.Decimal
	}{
		"it should be 385.40": {
			expectedValue: decimal.NewFromFloat(385.40),
		},
	}

	storage := NewMemoryStorage(context.Background())
	for i := range itineraries {
		storage.AddItinerary(itineraries[i])
	}

	for message, item := range items {
		itinerary, _ := storage.GetMostExpensive(source, destination)
		cmp := itinerary.GetPrice(entities.ChargeTypeTotalAmount, entities.TypeSingleAdult).Cmp(item.expectedValue)
		assert.Equal(t, decimalEqualNum, cmp, message)
	}
}

func TestService_GetLongest(t *testing.T) {
	items := map[string]struct {
		expectedValue int64
	}{
		"it should be 19 hours": {
			expectedValue: int64(19 * time.Hour),
		},
	}

	storage := NewMemoryStorage(context.Background())
	for i := range itineraries {
		storage.AddItinerary(itineraries[i])
	}

	for message, item := range items {
		itinerary, _ := storage.GetLongest(source, destination)
		assert.Equal(t, item.expectedValue, itinerary.GetDuration(), message)
	}
}

func TestService_GetShortest(t *testing.T) {
	items := map[string]struct {
		expectedValue int64
	}{
		"it should be 16 hours": {
			expectedValue: int64(16 * time.Hour),
		},
	}

	storage := NewMemoryStorage(context.Background())
	for i := range itineraries {
		storage.AddItinerary(itineraries[i])
	}

	for message, item := range items {
		itinerary, _ := storage.GetShortest(source, destination)
		assert.Equal(t, item.expectedValue, itinerary.GetDuration(), message)
	}
}
