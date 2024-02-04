package tripService

import (
	"fmt"
	"math"
	"trip/internals/database"
	"trip/internals/tripDetails"
	"trip/types"
)

func ListTrips(dto types.ListTripsParamsDto) ([]types.Trip, error) {
	trips, err := database.ListTrips(dto)

	if err != nil {
		fmt.Println("Error reading the response:", err)
		return trips, err
	}
	return trips, nil
}

func CreateTrip(dto types.CreateTripDto) (types.Trip, error) {
	price, calcPriceErr := CalcTripPrice(types.TripPriceDto{
		Origin:      dto.Origin,
		Destination: dto.Destination,
	})
	if calcPriceErr != nil {
		fmt.Println("Error reading the response:", calcPriceErr)
		return types.Trip{}, calcPriceErr
	}
	trip, createTripErr := database.CreateTrip(types.Trip{
		Origin:      dto.Origin,
		Destination: dto.Destination,
		Status:      `pending`,
		Price:       price,
	})
	if createTripErr != nil {
		fmt.Println("Got error while creating new trip item:", createTripErr)
		return types.Trip{}, createTripErr
	}
	return trip, nil
}

func CalcTripPrice(dto types.TripPriceDto) (float64, error) {
	pricePerKmUnit := 0.35
	pricePerMinUnit := 0.1
	const meterToKm = float64(1) / float64(1000)
	const secToMin = float64(1) / float64(60)
	tripData, err := tripDetails.GetTripInfo(
		types.Location{
			Origin:      dto.Origin,
			Destination: dto.Destination},
	)
	if err != nil {
		fmt.Println("Error Calculating the price:", err)
		return 0, err
	}
	pricePerDistance := pricePerKmUnit * (float64(tripData.Distance.Value) * meterToKm)
	pricePerTime := pricePerMinUnit * (float64(tripData.Duration.Value) * secToMin)
	price := pricePerDistance + pricePerTime
	return math.Round(price*100) / 100, nil
}
