package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	. "github.com/webapps/ataxi"

	"github.com/kellydunn/golang-geo"
)

func handlePassenger(taxis []*Taxi, potentialTaxis []*Taxi, passenger *Passenger, maxOccupancy int32) ([]*Taxi, []*Taxi) {
	availableTaxis := []*Taxi{}
	for _, taxi := range potentialTaxis {
		if !taxi.HasDeparted(passenger.DepartureTime) && !taxi.IsFull() {
			availableTaxis = append(availableTaxis, taxi)
		} else if taxi.NumPassengers() == 1 {
			taxi.DepartureTime = passenger.DepartureTime
		}
	}
	potentialTaxis = availableTaxis
	taxi, newTaxiStand := passenger.AssignedTaxi(potentialTaxis)
	if taxi == nil {
		taxi := &Taxi{
			ID:            int32(len(taxis)),
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
	taxis := []*Taxi{}
	potentialTaxis := []*Taxi{}

	start := time.Now()

	csvFile, _ := os.Open("../data/san_mateo_trips.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	i := 0
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		row := ParseLine(line)
		oCoord := Coord{X: row.OXCoord, Y: row.OYCoord}
		dCoord := Coord{X: row.DXCoord, Y: row.DYCoord}
		oLatLon := geo.NewPoint(row.OLat, row.OLon)
		dLatLon := geo.NewPoint(row.DLat, row.DLon)
		tripDistance := GetTripDistance(oLatLon, dLatLon)
		tripCategory := GetTripCategory(tripDistance)
		dSuperPixel := GetSuperPixel(dCoord, tripCategory)
		passenger := NewPassenger(
			row.PersonID,
			row.OFIPS,
			oCoord,
			oLatLon,
			row.DFIPS,
			dCoord,
			dLatLon,
			row.ODepartureTime,
			tripCategory,
			tripDistance,
			dSuperPixel,
		)
		if tripCategory == 0 {
			continue
		}
		taxis, potentialTaxis = handlePassenger(taxis, potentialTaxis, passenger, 5)
		fmt.Printf("\r%d - %d - %d", i, len(taxis), len(potentialTaxis))
		i++
	}
	fmt.Println()
	var pm float64
	var vm float64
	var numPassengers int32
	for _, taxi := range taxis {
		/*
			for _, passenger := range taxi.Passengers {
				if passenger.OCoord != taxi.OCoord {
					panic("passenger has different OCoord")
				}
				if passenger.DSuperPixel != taxi.DSuperPixel {
					panic("passenger has different OCoord")
				}
				if passenger.DepartureTime > taxi.DepartureTime {
					panic("taxi departs before passenger departure time")
				}
			}*/
		pm += taxi.PassengerTripMiles()
		vm += taxi.VehicleTripMiles()
		numPassengers += taxi.NumPassengers()
		//fmt.Println(taxi)
		/*
			for j := i - 1; j >= 0; j-- {
				taxi2 := taxis[j]
				if taxi.OCoord != taxi2.OCoord {
					break
				}
				if taxi.DSuperPixel == taxi2.DSuperPixel && !taxi2.IsFull() {
					for _, passenger := range taxi.Passengers {
						if !taxi2.HasDeparted(passenger.DepartureTime) &&
							taxi2.HasDeparted(passenger.LatestPickUpTime) {
							fmt.Println(taxi2)
							fmt.Println(taxi)
							panic("existing taxi was not full but new taxi was allocated")
						}
					}
				}
			}*/
	}
	fmt.Println(numPassengers)
	fmt.Println(float64(numPassengers) / float64(len(taxis)))
	fmt.Printf("AVO: %f\n", pm/vm)
	elapsed := time.Since(start)
	fmt.Printf("AVO analysis took %s\n", elapsed)
}
