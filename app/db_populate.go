package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	. "github.com/webapps/ataxi"
)

const bufferSize = 10000
const batchSize = 100
const numGoRoutines = 100

func handlePassenger(db *gorm.DB, taxis []*Taxi, potentialTaxis []*Taxi, passenger *Passenger, maxOccupancy uint32) ([]*Taxi, []*Taxi) {
	availableTaxis := []*Taxi{}
	departedTaxis := []*Taxi{}
	for _, taxi := range potentialTaxis {
		if !taxi.HasDeparted(passenger.DepartureTime) && !taxi.IsFull() {
			availableTaxis = append(availableTaxis, taxi)
		} else {
			if taxi.NumPassengers() == 1 {
				taxi.DepartureTime = passenger.DepartureTime
			}
			departedTaxis = append(departedTaxis, taxi)
		}
	}

	if db != nil {
		var wg sync.WaitGroup
		createTaxis(db, departedTaxis, &wg)
		wg.Wait()
	}

	potentialTaxis = availableTaxis
	taxi, newTaxiStand := passenger.AssignedTaxi(potentialTaxis)
	if taxi == nil {
		taxi := NewTaxi(passenger, maxOccupancy)
		taxis = append(taxis, taxi)
		if newTaxiStand {
			if db != nil {
				var wg sync.WaitGroup
				createTaxis(db, potentialTaxis, &wg)
				wg.Wait()
			}
			potentialTaxis = []*Taxi{taxi}
		} else {
			potentialTaxis = append(potentialTaxis, taxi)
		}
	} else {
		taxi.Passengers = append(taxi.Passengers, passenger)
	}
	return taxis, potentialTaxis
}

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func createPassengers(db *gorm.DB, passengers []*Passenger, wg *sync.WaitGroup) {
	start, length := 0, len(passengers)
	var end int
	for i := 0; i < min(numGoRoutines, len(passengers)/batchSize); i++ {
		end = min(length, start+batchSize)
		wg.Add(1)
		go func(batch []*Passenger) {
			defer wg.Done()
			for _, passenger := range batch {
				db.Create(passenger)
			}
		}(passengers[start:end])
		start += batchSize
	}
}

func createTaxis(db *gorm.DB, taxis []*Taxi, wg *sync.WaitGroup) {
	for _, taxi := range taxis {
		wg.Add(1)
		go func(t *Taxi) {
			defer wg.Done()
			db.Create(t)
		}(taxi)
	}
	/*
		start, length := 0, len(taxis)
		var end int
		for i := 0; i < min(numGoRoutines, len(taxis)/batchSize); i++ {
			end = min(length, start+batchSize)
			wg.Add(1)
			go func(batch []*Taxi) {
				defer wg.Done()
				for _, taxi := range batch {
					db.Create(taxi)
				}
			}(taxis[start:end])
			start += batchSize
		}*/
}

func main() {
	start := time.Now()

	raw, err := ioutil.ReadFile("../config.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	var config Config
	json.Unmarshal(raw, &config)

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8&parseTime=True&loc=Local", config.Username, config.Password, config.Database))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer db.Close()

	db.DropTable(&Passenger{}, &Taxi{})
	db.AutoMigrate(&Passenger{}, &Taxi{})

	csvFile, _ := os.Open("../data/san_mateo_trips.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	fmt.Println("Reading trip csv...")
	var taxis []*Taxi
	var potentialTaxis []*Taxi
	var passengers []*Passenger
	var processed int
	for {
		line, err := reader.Read()
		if err == io.EOF {
			if len(passengers) > 0 {
				var wg sync.WaitGroup
				createPassengers(db, passengers, &wg)
				wg.Wait()
				processed += len(passengers)
				fmt.Printf("\rProcessed %d", processed)
				passengers = []*Passenger{}
			}
			break
		} else if err != nil {
			log.Fatal(err)
		}
		row := ParseLine(line)
		passenger := NewPassengerFromRow(row)
		if passenger.TripCategory == 0 {
			continue
		}
		passengers = append(passengers, passenger)

		if len(passengers) == bufferSize {
			var wg sync.WaitGroup
			createPassengers(db, passengers, &wg)
			wg.Wait()
			processed += len(passengers)
			fmt.Printf("\rProcessed %d", processed)
			passengers = []*Passenger{}
		}

		taxis, potentialTaxis = handlePassenger(db, taxis, potentialTaxis, passenger, 5)
	}

	if len(potentialTaxis) > 0 {
		var wg sync.WaitGroup
		createTaxis(db, potentialTaxis, &wg)
		wg.Wait()
	}

	fmt.Println()
	var pm float64
	var vm float64
	var numPassengers uint32
	fmt.Printf("Num of Taxis: %d\n", len(taxis))
	for i, taxi := range taxis {
		pm += taxi.PassengerTripMiles()
		vm += taxi.VehicleTripMiles()
		numPassengers += taxi.NumPassengers()
		fmt.Printf("\rProcessed %d taxi(s)", i)

		// for j := i - 1; j >= 0; j-- {
		// 	taxi2 := taxis[j]
		// 	if taxi.OX != taxi2.OX || taxi.OY != taxi2.OY {
		// 		break
		// 	}
		// 	if taxi.DXSuper == taxi2.DXSuper &&
		// 		taxi.DYSuper == taxi2.DYSuper &&
		// 		!taxi2.IsFull() {
		// 		for _, passenger := range taxi.Passengers {
		// 			if !taxi2.HasDeparted(passenger.DepartureTime) &&
		// 				taxi2.HasDeparted(passenger.LatestPickUpTime) {
		// 				fmt.Println(taxi2)
		// 				fmt.Println(taxi)
		// 				panic("existing taxi was not full but new taxi was allocated")
		// 			}
		// 		}
		// 	}
		// }
	}
	fmt.Println()
	fmt.Println(numPassengers)
	fmt.Println(float64(numPassengers) / float64(len(taxis)))
	fmt.Printf("AVO: %f\n", pm/vm)
	elapsed := time.Since(start)
	fmt.Printf("AVO analysis took %s\n", elapsed)
}
