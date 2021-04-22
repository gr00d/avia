package entities

import (
	"encoding/xml"
	"time"

	"github.com/shopspring/decimal"
)

type SourceCity string
type DestinationCity string
type ItineraryUUID string

type Itineraries struct {
	Itineraries   []*Itinerary
	Shortest      *Itinerary
	Longest       *Itinerary
	Cheapest      *Itinerary
	MostExpensive *Itinerary
	Optimal       *Itinerary
}

type SearchResponse struct {
	RequestID   string      `xml:"RequestId"`
	Itineraries []Itinerary `xml:"PricedItineraries>Flights"`
}

type Itinerary struct {
	UUID    ItineraryUUID
	Onward  []Flight `xml:"OnwardPricedItinerary>Flights>Flight"`
	Return  []Flight `xml:"ReturnPricedItinerary>Flights>Flight"`
	Pricing *Price   `xml:"Pricing"`
}

type Flight struct {
	Carrier            string
	FlightNumber       string
	Source             string
	Destination        string
	DepartureTimeStamp FlightDate
	ArrivalTimeStamp   FlightDate
	Class              string
	NumberOfStops      string
	TicketType         string
}

type FlightDate struct {
	time.Time
}

type Price struct {
	Currency       string   `xml:"currency,attr"`
	ServiceCharges []Charge `xml:"ServiceCharges"`
}

type Charge struct {
	ChargeType string          `xml:"ChargeType,attr"`
	Type       string          `xml:"type,attr"`
	Cost       decimal.Decimal `xml:",chardata"`
}

const (
	ChargeTypeBaseFare     = "BaseFare"
	ChargeTypeAirlineTaxes = "AirlineTaxes"
	ChargeTypeTotalAmount  = "TotalAmount"
)

const (
	TypeSingleAdult  = "SingleAdult"
	TypeSingleChild  = "SingleChild"
	TypeSingleInfant = "SingleInfant"
)

func (c *FlightDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	_ = d.DecodeElement(&v, &start)
	parse, err := time.Parse("2006-01-02T1504", v)
	if err != nil {
		return nil
	}
	*c = FlightDate{parse}
	return nil
}

func (i *Itinerary) GetPrice(chargeType, rateType string) decimal.Decimal {
	if i.Pricing != nil {
		for s := range i.Pricing.ServiceCharges {
			if i.Pricing.ServiceCharges[s].ChargeType == chargeType && i.Pricing.ServiceCharges[s].Type == rateType {
				return i.Pricing.ServiceCharges[s].Cost
			}
		}
	}

	return decimal.Zero
}

func (i *Itinerary) GetDuration() int64 {
	if len(i.Onward) == 0 {
		return 0
	}
	departureTime, arrivalTime := i.Onward[0].DepartureTimeStamp.Time, i.Onward[len(i.Onward)-1].ArrivalTimeStamp.Time
	return int64(arrivalTime.Sub(departureTime))
}

func (i *Itinerary) GetDurationWithoutTransfer() int64 {
	if len(i.Onward) == 0 {
		return 0
	}

	var result int64
	for k := range i.Onward {
		departureTime, arrivalTime := i.Onward[k].DepartureTimeStamp.Time, i.Onward[k].ArrivalTimeStamp.Time
		result += int64(arrivalTime.Sub(departureTime))
	}
	return result
}

func (i *Itinerary) GetTransferDuration() int64 {
	if len(i.Onward) == 0 {
		return 0
	}

	var result int64
	for k := 0; k < len(i.Onward)-1; k++ {
		departureTime, arrivalTime := i.Onward[k+1].DepartureTimeStamp.Time, i.Onward[k].ArrivalTimeStamp.Time
		result += int64(departureTime.Sub(arrivalTime))
	}
	return result
}
