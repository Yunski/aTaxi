package ataxi

import (
	"fmt"
	"math"

	geo "github.com/kellydunn/golang-geo"
)

type Taxi struct {
	ID            int32
	OCoord        Coord
	OLatLon       *geo.Point
	DCoord        Coord
	DLatLon       *geo.Point
	Passengers    []*Passenger
	DepartureTime int32
	MaxOccupancy  int32
	DSuperPixel   Coord
	NaiveVM       bool
}

func (taxi *Taxi) HasDeparted(time int32) bool {
	return taxi.DepartureTime <= time
}

func (taxi *Taxi) NumPassengers() int32 {
	return int32(len(taxi.Passengers))
}

func (taxi *Taxi) IsFull() bool {
	return taxi.NumPassengers() == taxi.MaxOccupancy
}

func (taxi *Taxi) VehicleTripMiles() float64 {
	curLatLon := taxi.OLatLon
	passengers := taxi.Passengers
	vtm := 0.0
	if taxi.NaiveVM {
		for _, passenger := range passengers {
			vtm += passenger.DistanceTo(curLatLon)
			curLatLon = passenger.DLatLon
		}
	} else {
		for len(passengers) > 0 {
			nearestDist := math.Inf(1)
			nearestIdx, idx := 0, 0
			for _, passenger := range passengers {
				dist := passenger.DistanceTo(curLatLon)
				if dist < nearestDist {
					nearestDist = dist
					nearestIdx = idx
				}
				idx++
			}
			vtm += nearestDist
			curLatLon = passengers[nearestIdx].DLatLon
			remaining := []*Passenger{}
			for i, passenger := range passengers {
				if i != nearestIdx {
					remaining = append(remaining, passenger)
				}
			}
			passengers = remaining
		}
	}
	return vtm
}

func (taxi *Taxi) PassengerTripMiles() float64 {
	ptm := 0.0
	for _, passenger := range taxi.Passengers {
		ptm += passenger.TripDistance
	}
	return ptm
}

func (taxi *Taxi) String() string {
	return fmt.Sprintf("ID: %d DepartureTime: %d NumPassengers: %d VO: %f Passenger 1 DepartureTime: %d, OCoord: %s, DSuper: %s",
		taxi.ID,
		taxi.DepartureTime,
		taxi.NumPassengers(),
		taxi.PassengerTripMiles()/taxi.VehicleTripMiles(),
		taxi.Passengers[0].DepartureTime,
		taxi.OCoord,
		taxi.DSuperPixel,
	)
}
