package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/webapps/ataxi"
)

func handlePassenger(db *gorm.DB, taxis []*ataxi.Taxi, potentialTaxis []*ataxi.Taxi,
	passenger *ataxi.Passenger, maxOccupancy uint32) ([]*ataxi.Taxi, []*ataxi.Taxi) {
	var availableTaxis []*ataxi.Taxi
	for _, taxi := range potentialTaxis {
		if !taxi.HasDeparted(passenger.DepartureTime) && !taxi.IsFull() {
			availableTaxis = append(availableTaxis, taxi)
		} else {
			if taxi.NumPassengers == 1 {
				taxi.DepartureTime = passenger.DepartureTime
			}
			taxi.PMT = taxi.PersonMilesTraveled()
			taxi.VMT = taxi.VehicleMilesTraveled()
		}
	}

	potentialTaxis = availableTaxis
	taxi, newTaxiStand := passenger.FindTaxi(potentialTaxis)
	if taxi == nil {
		taxi = ataxi.NewTaxi(uint(len(taxis)+1), passenger, maxOccupancy)
		taxis = append(taxis, taxi)
		if newTaxiStand {
            for _, taxi := range potentialTaxis {
                taxi.PMT = taxi.PersonMilesTraveled()
                taxi.VMT = taxi.VehicleMilesTraveled()
            }
			potentialTaxis = []*ataxi.Taxi{taxi}
		} else {
			potentialTaxis = append(potentialTaxis, taxi)
		}
	} else {
		taxi.AddPassenger(passenger)
	}
	return taxis, potentialTaxis
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(errors.New("You must provide a csv file."))
		os.Exit(1)
	}

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8&parseTime=True&loc=Local", ataxi.Config.Username, ataxi.Config.Password, ataxi.Config.Database))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer db.Close()

	if !db.HasTable(&ataxi.Passenger{}) && !db.HasTable(&ataxi.Taxi{}) {
		db.AutoMigrate(&ataxi.Passenger{}, &ataxi.Taxi{})
	}

	csvFileName := os.Args[1]
	csvFile, _ := os.Open(fmt.Sprintf("../data/%s", csvFileName))
	reader := csv.NewReader(bufio.NewReader(csvFile))

	start := time.Now()
	fmt.Println("Reading trip csv...")
	var taxis []*ataxi.Taxi
	var potentialTaxis []*ataxi.Taxi
	var id uint
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		row := ataxi.ParseLine(line)
		passenger := ataxi.NewPassengerFromRow(id+1, row)
		if passenger.TripCategory == 0 {
			continue
		}
		id++

		taxis, potentialTaxis = handlePassenger(db, taxis, potentialTaxis, passenger, 5)
	}
    for _, taxi := range potentialTaxis {
        taxi.PMT = taxi.PersonMilesTraveled()
        taxi.VMT = taxi.VehicleMilesTraveled()
    }

	fmt.Printf("Num of Taxis needed: %d\n", len(taxis))
	elapsed := time.Since(start)
	fmt.Printf("csv processing took %s\n", elapsed)
	start = time.Now()
	var pmt float64
	var vmt float64
	var numPassengers uint32
	for i, taxi := range taxis {
		db.Create(taxi)
		pmt += taxi.PMT
		vmt += taxi.VMT
		numPassengers += taxi.NumPassengers
		if i%10000 == 0 {
			fmt.Printf("\rProcessed %d taxi(s)", i)
		}
	}
	fmt.Println()
	fmt.Printf("\rFinished processing %d taxi(s)\n", len(taxis))
	fmt.Printf("Capacity ratio: %f\n", float64(numPassengers)/float64(len(taxis)))
	fmt.Printf("AVO: %f\n", pmt/vmt)
	elapsed = time.Since(start)
	fmt.Printf("AVO analysis took %s\n", elapsed)
}
