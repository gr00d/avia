package entities

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestItinerary_GetDuration(t *testing.T) {
	items := map[string]struct {
		xml           string
		expectedValue int64
	}{
		"it should be 16 hours with 2 flight": {
			xml: `
<Flights>
	<OnwardPricedItinerary>
		<Flights>
			<Flight>
				<Carrier id="CZ">China Southern Airlines</Carrier>
				<FlightNumber>384</FlightNumber>
				<Source>DXB</Source>
				<Destination>CAN</Destination>
				<DepartureTimeStamp>2018-10-27T0100</DepartureTimeStamp>
				<ArrivalTimeStamp>2018-10-27T1200</ArrivalTimeStamp>
				<Class>T</Class>
				<NumberOfStops>0</NumberOfStops>
				<FareBasis>
				2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_CAN_384_107_01:40_$255_CAN_BKK_363_107_16:00__A2_1_1
				</FareBasis>
				<WarningText/>
				<TicketType>E</TicketType>
			</Flight>
			<Flight>
				<Carrier id="CZ">China Southern Airlines</Carrier>
				<FlightNumber>363</FlightNumber>
				<Source>CAN</Source>
				<Destination>BKK</Destination>
				<DepartureTimeStamp>2018-10-27T1600</DepartureTimeStamp>
				<ArrivalTimeStamp>2018-10-27T1700</ArrivalTimeStamp>
				<Class>Y</Class>
				<NumberOfStops>0</NumberOfStops>
				<FareBasis>
				2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_CAN_384_107_01:40_$255_CAN_BKK_363_107_16:00__A2_1_1
				</FareBasis>
				<WarningText/>
				<TicketType>E</TicketType>
			</Flight>
		</Flights>
	</OnwardPricedItinerary>
	<Pricing currency="SGD">
		<ServiceCharges type="SingleAdult" ChargeType="BaseFare">233.00</ServiceCharges>
	</Pricing>
</Flights>`,
			expectedValue: int64(16 * time.Hour),
		},
		"it should be 1 hour with 1 light": {
			xml: `
<Flights>
	<OnwardPricedItinerary>
		<Flights>
			<Flight>
				<Carrier id="CZ">China Southern Airlines</Carrier>
				<FlightNumber>384</FlightNumber>
				<Source>DXB</Source>
				<Destination>CAN</Destination>
				<DepartureTimeStamp>2018-10-28T0100</DepartureTimeStamp>
				<ArrivalTimeStamp>2018-10-28T0200</ArrivalTimeStamp>
				<Class>T</Class>
				<NumberOfStops>0</NumberOfStops>
				<FareBasis>
				2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_CAN_384_107_01:40_$255_CAN_BKK_363_107_16:00__A2_1_1
				</FareBasis>
				<WarningText/>
				<TicketType>E</TicketType>
			</Flight>
		</Flights>
	</OnwardPricedItinerary>
	<Pricing currency="SGD">
		<ServiceCharges type="SingleAdult" ChargeType="BaseFare">233.00</ServiceCharges>
	</Pricing>
</Flights>`,
			expectedValue: int64(1 * time.Hour),
		},
		"it should be 0 hour with 0 light": {
			xml: `
<Flights>
	<OnwardPricedItinerary>
		<Flights>
		</Flights>
	</OnwardPricedItinerary>
	<Pricing currency="SGD">
		<ServiceCharges type="SingleAdult" ChargeType="BaseFare">233.00</ServiceCharges>
	</Pricing>
</Flights>`,
			expectedValue: int64(0 * time.Hour),
		},
	}

	for message, item := range items {
		itinerary := &Itinerary{}
		_ = xml.Unmarshal([]byte(item.xml), itinerary)
		assert.Equal(t, item.expectedValue, itinerary.GetDuration(), message)
	}
}

func TestItinerary_GetDurationWithoutTransfer(t *testing.T) {
	items := map[string]struct {
		xml           string
		expectedValue int64
	}{
		"it should be 12 hours with 2 flight": {
			xml: `
<Flights>
	<OnwardPricedItinerary>
		<Flights>
			<Flight>
				<Carrier id="CZ">China Southern Airlines</Carrier>
				<FlightNumber>384</FlightNumber>
				<Source>DXB</Source>
				<Destination>CAN</Destination>
				<DepartureTimeStamp>2018-10-27T0100</DepartureTimeStamp>
				<ArrivalTimeStamp>2018-10-27T1200</ArrivalTimeStamp>
				<Class>T</Class>
				<NumberOfStops>0</NumberOfStops>
				<FareBasis>
				2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_CAN_384_107_01:40_$255_CAN_BKK_363_107_16:00__A2_1_1
				</FareBasis>
				<WarningText/>
				<TicketType>E</TicketType>
			</Flight>
			<Flight>
				<Carrier id="CZ">China Southern Airlines</Carrier>
				<FlightNumber>363</FlightNumber>
				<Source>CAN</Source>
				<Destination>BKK</Destination>
				<DepartureTimeStamp>2018-10-27T1600</DepartureTimeStamp>
				<ArrivalTimeStamp>2018-10-27T1700</ArrivalTimeStamp>
				<Class>Y</Class>
				<NumberOfStops>0</NumberOfStops>
				<FareBasis>
				2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_CAN_384_107_01:40_$255_CAN_BKK_363_107_16:00__A2_1_1
				</FareBasis>
				<WarningText/>
				<TicketType>E</TicketType>
			</Flight>
		</Flights>
	</OnwardPricedItinerary>
	<Pricing currency="SGD">
		<ServiceCharges type="SingleAdult" ChargeType="BaseFare">233.00</ServiceCharges>
	</Pricing>
</Flights>`,
			expectedValue: int64(12 * time.Hour),
		},
		"it should be 1 hour with 1 light": {
			xml: `
<Flights>
	<OnwardPricedItinerary>
		<Flights>
			<Flight>
				<Carrier id="CZ">China Southern Airlines</Carrier>
				<FlightNumber>384</FlightNumber>
				<Source>DXB</Source>
				<Destination>CAN</Destination>
				<DepartureTimeStamp>2018-10-28T0100</DepartureTimeStamp>
				<ArrivalTimeStamp>2018-10-28T0200</ArrivalTimeStamp>
				<Class>T</Class>
				<NumberOfStops>0</NumberOfStops>
				<FareBasis>
				2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_CAN_384_107_01:40_$255_CAN_BKK_363_107_16:00__A2_1_1
				</FareBasis>
				<WarningText/>
				<TicketType>E</TicketType>
			</Flight>
		</Flights>
	</OnwardPricedItinerary>
	<Pricing currency="SGD">
		<ServiceCharges type="SingleAdult" ChargeType="BaseFare">233.00</ServiceCharges>
	</Pricing>
</Flights>`,
			expectedValue: int64(1 * time.Hour),
		},
		"it should be 0 hour with 0 light": {
			xml: `
<Flights>
	<OnwardPricedItinerary>
		<Flights>
		</Flights>
	</OnwardPricedItinerary>
	<Pricing currency="SGD">
		<ServiceCharges type="SingleAdult" ChargeType="BaseFare">233.00</ServiceCharges>
	</Pricing>
</Flights>`,
			expectedValue: int64(0 * time.Hour),
		},
	}

	for message, item := range items {
		itinerary := &Itinerary{}
		_ = xml.Unmarshal([]byte(item.xml), itinerary)
		assert.Equal(t, item.expectedValue, itinerary.GetDurationWithoutTransfer(), message)
	}
}

func TestItinerary_GetTransferDuration(t *testing.T) {
	items := map[string]struct {
		xml           string
		expectedValue int64
	}{
		"it should be 4 hours with 2 flight": {
			xml: `
<Flights>
	<OnwardPricedItinerary>
		<Flights>
			<Flight>
				<Carrier id="CZ">China Southern Airlines</Carrier>
				<FlightNumber>384</FlightNumber>
				<Source>DXB</Source>
				<Destination>CAN</Destination>
				<DepartureTimeStamp>2018-10-27T0100</DepartureTimeStamp>
				<ArrivalTimeStamp>2018-10-27T1200</ArrivalTimeStamp>
				<Class>T</Class>
				<NumberOfStops>0</NumberOfStops>
				<FareBasis>
				2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_CAN_384_107_01:40_$255_CAN_BKK_363_107_16:00__A2_1_1
				</FareBasis>
				<WarningText/>
				<TicketType>E</TicketType>
			</Flight>
			<Flight>
				<Carrier id="CZ">China Southern Airlines</Carrier>
				<FlightNumber>363</FlightNumber>
				<Source>CAN</Source>
				<Destination>BKK</Destination>
				<DepartureTimeStamp>2018-10-27T1600</DepartureTimeStamp>
				<ArrivalTimeStamp>2018-10-27T1700</ArrivalTimeStamp>
				<Class>Y</Class>
				<NumberOfStops>0</NumberOfStops>
				<FareBasis>
				2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_CAN_384_107_01:40_$255_CAN_BKK_363_107_16:00__A2_1_1
				</FareBasis>
				<WarningText/>
				<TicketType>E</TicketType>
			</Flight>
		</Flights>
	</OnwardPricedItinerary>
	<Pricing currency="SGD">
		<ServiceCharges type="SingleAdult" ChargeType="BaseFare">233.00</ServiceCharges>
	</Pricing>
</Flights>`,
			expectedValue: int64(4 * time.Hour),
		},
		"it should be 0 hour with 1 light": {
			xml: `
<Flights>
	<OnwardPricedItinerary>
		<Flights>
			<Flight>
				<Carrier id="CZ">China Southern Airlines</Carrier>
				<FlightNumber>384</FlightNumber>
				<Source>DXB</Source>
				<Destination>CAN</Destination>
				<DepartureTimeStamp>2018-10-28T0100</DepartureTimeStamp>
				<ArrivalTimeStamp>2018-10-28T0200</ArrivalTimeStamp>
				<Class>T</Class>
				<NumberOfStops>0</NumberOfStops>
				<FareBasis>
				2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_CAN_384_107_01:40_$255_CAN_BKK_363_107_16:00__A2_1_1
				</FareBasis>
				<WarningText/>
				<TicketType>E</TicketType>
			</Flight>
		</Flights>
	</OnwardPricedItinerary>
	<Pricing currency="SGD">
		<ServiceCharges type="SingleAdult" ChargeType="BaseFare">233.00</ServiceCharges>
	</Pricing>
</Flights>`,
			expectedValue: int64(0 * time.Hour),
		},
		"it should be 0 hour with 0 light": {
			xml: `
<Flights>
	<OnwardPricedItinerary>
		<Flights>
		</Flights>
	</OnwardPricedItinerary>
	<Pricing currency="SGD">
		<ServiceCharges type="SingleAdult" ChargeType="BaseFare">233.00</ServiceCharges>
	</Pricing>
</Flights>`,
			expectedValue: int64(0 * time.Hour),
		},
	}

	for message, item := range items {
		itinerary := &Itinerary{}
		_ = xml.Unmarshal([]byte(item.xml), itinerary)
		assert.Equal(t, item.expectedValue, itinerary.GetTransferDuration(), message)
	}
}
