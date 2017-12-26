package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/kellydunn/golang-geo"
)

type Passenger struct {
	PersonID           int
	OFIPS              int
	OCoord             Coord
	OLatLon            *geo.Point
	DFIPS              int
	DCoord             Coord
	DLatLon            *geo.Point
	DepartureTime      int
	TripDistance       float64
	TripCategory       int
	DistanceConstraint float64
	LatestPickUpTime   int
	DSuperPixel        Coord
}

func (passenger *Passenger) AssignedTaxi(taxis []*Taxi) (*Taxi, bool) {
	if len(taxis) == 0 || passenger.OCoord != taxis[0].OCoord {
		return nil, true
	}
	//constraint := passenger.DistanceConstraint
	var matchedTaxi *Taxi
	for _, taxi := range taxis {
		//if math.Abs(float64(passenger.DCoord.X-taxi.DCoord.X)) <= constraint &&
		//	math.Abs(float64(passenger.DCoord.Y-taxi.DCoord.Y)) <= constraint &&
		if taxi.DSuperPixel == passenger.DSuperPixel && taxi.HasDeparted(passenger.LatestPickUpTime) {
			if matchedTaxi == nil || taxi.DepartureTime < matchedTaxi.DepartureTime {
				matchedTaxi = taxi
			}
		}
	}
	return matchedTaxi, false
}

func (passenger *Passenger) DistanceTo(latlon *geo.Point) float64 {
	return getTripDistance(passenger.DLatLon, latlon)
}

type Taxi struct {
	ID            int
	OCoord        Coord
	OLatLon       *geo.Point
	DCoord        Coord
	DLatLon       *geo.Point
	Passengers    []*Passenger
	DepartureTime int
	MaxOccupancy  int
	DSuperPixel   Coord
}

func (taxi *Taxi) HasDeparted(time int) bool {
	return taxi.DepartureTime <= time
}

func (taxi *Taxi) NumPassengers() int {
	return len(taxi.Passengers)
}

func (taxi *Taxi) IsFull() bool {
	return taxi.NumPassengers() == taxi.MaxOccupancy
}

func (taxi *Taxi) VehicleTripMiles() float64 {
	curLatLon := taxi.OLatLon
	passengers := taxi.Passengers
	vtm := 0.0
	for _, passenger := range passengers {
		vtm += passenger.DistanceTo(curLatLon)
		curLatLon = passenger.DLatLon
	}
	/*
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
		}*/
	return vtm
}

func (taxi *Taxi) PassengerTripMiles() float64 {
	ptm := 0.0
	for _, passenger := range taxi.Passengers {
		ptm += passenger.TripDistance
	}
	return ptm
}

func handlePassenger(taxis []*Taxi, potentialTaxis []*Taxi, passenger *Passenger, maxOccupancy int) ([]*Taxi, []*Taxi) {
	availableTaxis := []*Taxi{}
	for _, taxi := range potentialTaxis {
		if !taxi.HasDeparted(passenger.DepartureTime) && !taxi.IsFull() {
			availableTaxis = append(availableTaxis, taxi)
		}
	}
	potentialTaxis = availableTaxis
	taxi, newTaxiStand := passenger.AssignedTaxi(potentialTaxis)
	if taxi == nil {
		taxi := &Taxi{
			ID:            len(taxis),
			OCoord:        passenger.OCoord,
			OLatLon:       passenger.OLatLon,
			DCoord:        passenger.DCoord,
			DLatLon:       passenger.DLatLon,
			Passengers:    []*Passenger{passenger},
			DepartureTime: passenger.LatestPickUpTime,
			MaxOccupancy:  maxOccupancy,
			DSuperPixel:   passenger.DSuperPixel,
		}
		taxis = append(taxis, taxi)
		if newTaxiStand {
			potentialTaxis = []*Taxi{taxi}
		} else {
			potentialTaxis = append(potentialTaxis, taxi)
		}

	} else {
		taxi.Passengers = append(taxi.Passengers, passenger)
	}
	return taxis, potentialTaxis
}

func main() {
	OFIPS := 5
	DFIPS := 13
	PersonID := 1
	OXCoord := 8
	OYCoord := 9
	OLon := 6
	OLat := 7
	DXCoord := 16
	DYCoord := 17
	DLon := 14
	DLat := 15
	ODepartureTime := 10

	MAX_OCCUPANCY := 5

	taxis := []*Taxi{}
	potentialTaxis := []*Taxi{}

	start := time.Now()

	csvFile, _ := os.Open("data/san_mateo_trips.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	numPassengers := 0
	for {
		row, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		oFIPS, _ := strconv.ParseInt(row[OFIPS], 0, 0)
		dFIPS, _ := strconv.ParseInt(row[DFIPS], 0, 0)
		personID, _ := strconv.ParseInt(row[PersonID], 0, 0)
		oX, _ := strconv.ParseInt(row[OXCoord], 0, 0)
		oY, _ := strconv.ParseInt(row[OYCoord], 0, 0)
		oCoord := Coord{
			X: int(oX),
			Y: int(oY),
		}
		dX, _ := strconv.ParseInt(row[DXCoord], 0, 0)
		dY, _ := strconv.ParseInt(row[DYCoord], 0, 0)
		dCoord := Coord{
			X: int(dX),
			Y: int(dY),
		}
		oLat, _ := strconv.ParseFloat(row[OLat], 64)
		oLon, _ := strconv.ParseFloat(row[OLon], 64)
		oLatLon := geo.NewPoint(oLat, oLon)
		dLat, _ := strconv.ParseFloat(row[DLat], 64)
		dLon, _ := strconv.ParseFloat(row[DLon], 64)
		dLatLon := geo.NewPoint(dLat, dLon)
		departureTime, _ := strconv.ParseInt(row[ODepartureTime], 0, 0)
		tripDistance := getTripDistance(oLatLon, dLatLon)
		tripCategory := getTripCategory(tripDistance)
		distConstraint := getDistConstraint(tripDistance)
		dSuperPixel := getSuperPixel(dCoord, tripCategory)
		passenger := &Passenger{
			OFIPS:              int(oFIPS),
			PersonID:           int(personID),
			OCoord:             oCoord,
			OLatLon:            oLatLon,
			DFIPS:              int(dFIPS),
			DCoord:             dCoord,
			DLatLon:            dLatLon,
			DepartureTime:      int(departureTime),
			TripCategory:       tripCategory,
			TripDistance:       tripDistance,
			DistanceConstraint: distConstraint,
			LatestPickUpTime:   int(departureTime) + getMaxWaitingTime(tripDistance),
			DSuperPixel:        dSuperPixel,
		}
		if tripCategory == 0 {
			numPassengers++
			continue
		}
		taxis, potentialTaxis = handlePassenger(taxis, potentialTaxis, passenger, MAX_OCCUPANCY)
		fmt.Printf("\r%d - %d - %d", numPassengers, len(taxis), len(potentialTaxis))
		numPassengers++
	}
	fmt.Println()
	tvo := 0.0
	for _, taxi := range taxis {
		pm, vm := taxi.PassengerTripMiles(), taxi.VehicleTripMiles()
		vo := pm / vm
		//fmt.Printf("num passengers: %d - pm: %f - vm: %f - VO: %f\n", len(taxi.Passengers), pm, vm, vo)
		tvo += vo
	}
	fmt.Println(numPassengers / len(taxis))
	fmt.Printf("AVO: %f\n", tvo/float64(len(taxis)))
	elapsed := time.Since(start)
	log.Printf("AVO analysis took %s\n", elapsed)
}
