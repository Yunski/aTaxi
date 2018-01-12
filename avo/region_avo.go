package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/webapps/ataxi"
)

func handlePassenger(taxis []*ataxi.Taxi, potentialTaxis []*ataxi.Taxi,
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

func getMT(taxis []*ataxi.Taxi) (float64, float64) {
	var pmt float64
	var vmt float64
	for _, taxi := range taxis {
		pmt += taxi.PMT
		vmt += taxi.VMT
	}
	return pmt, vmt
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(errors.New("You must provide a data directory containing the ataxi mode trip files."))
		os.Exit(1)
	}

	files, err := filepath.Glob(filepath.Join(os.Args[1], "*.csv"))
	if err != nil {
		log.Fatal(err)
	}

	stateFile, err := os.Create("../data/state_avos.csv")
	if err != nil {
		log.Fatal(err)
	}
	stateWriter := csv.NewWriter(stateFile)
	stateColumns := []string{"State", "AVO", "PMT", "VMT"}
	stateWriter.Write(stateColumns)
	var stateRow [4]string

	countyFile, err := os.Create("../data/county_avos.csv")
	if err != nil {
		log.Fatal(err)
	}
	countyWriter := csv.NewWriter(countyFile)
	countyColumns := []string{"County", "AVO", "PMT", "VMT"}
	countyWriter.Write(countyColumns)
	var countyRow [4]string

	tripFile, err := os.Create("../data/ataxi_trips.csv")
	if err != nil {
		log.Fatal(err)
	}
	tripWriter := csv.NewWriter(tripFile)
	tripColumns := []string{"OX", "OY", "DepartureTime", "DX", "DY",
		"MadeEmptyTime", "VehicleTripMiles", "DepartureOccupancy",
		"OccupantTripMiles", "OXSuper5", "OYSuper5", "DXSuper5", "DYSuper5",
		"OXSuper10", "OYSuper10", "DXSuper10", "DYSuper10"}
	tripWriter.Write(tripColumns)
	var tripRow [17]string

	start := time.Now()
	fmt.Println("Reading mode ataxi trip files...")

	var regionPMT float64
	var regionVMT float64
	var statePMT float64
	var stateVMT float64

	var stateFIPS string
	re := regexp.MustCompile("[0-9]+")
	var id uint
	for i, file := range files {
		var countyTaxis []*ataxi.Taxi
		var potentialTaxis []*ataxi.Taxi

		_, filename := filepath.Split(file)
		fmt.Printf("Processing %s\n", filename)
		curFIPS := re.FindAllString(filename, 1)[0][:2]
		if stateFIPS != curFIPS && i != 0 {
			stateRow[0] = stateFIPS
			stateRow[1] = strconv.FormatFloat(statePMT/stateVMT, 'f', 2, 64)
            stateRow[2] = strconv.FormatFloat(statePMT, 'f', 2, 64)
            stateRow[3] = strconv.FormatFloat(stateVMT, 'f', 2, 64)
			stateWriter.Write(stateRow[:])
			fmt.Printf("state %s avo: %s - pmt: %s - vmt: %s\n", stateRow[0], stateRow[1], stateRow[2], stateRow[3])
			statePMT = 0
			stateVMT = 0
		}
		stateFIPS = curFIPS

		csvFile, _ := os.Open(file)
		reader := csv.NewReader(bufio.NewReader(csvFile))
		if _, err := reader.Read(); err != nil {
			log.Fatal(err)
		}
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

			countyTaxis, potentialTaxis = handlePassenger(countyTaxis, potentialTaxis, passenger, 5)
		}
	    for _, taxi := range potentialTaxis {
            taxi.PMT = taxi.PersonMilesTraveled()
			taxi.VMT = taxi.VehicleMilesTraveled()
		}
	    pmt, vmt := getMT(countyTaxis)
		countyRow[0] = strconv.Itoa(int(countyTaxis[0].OFIPS))
		countyRow[1] = strconv.FormatFloat(pmt/vmt, 'f', 2, 64)
        countyRow[2] = strconv.FormatFloat(pmt, 'f', 2, 64)
        countyRow[3] = strconv.FormatFloat(vmt, 'f', 2, 64)
		countyWriter.Write(countyRow[:])
		fmt.Printf("county avo: %s - pmt: %s - vmt: %s\n", countyRow[1], countyRow[2], countyRow[3])

		for _, taxi := range countyTaxis {
            tripRow[0] = strconv.Itoa(int(taxi.OX))
			tripRow[1] = strconv.Itoa(int(taxi.OY))
			tripRow[2] = strconv.Itoa(int(taxi.DepartureTime) % 86400)
			tripRow[3] = strconv.Itoa(int(taxi.DX))
			tripRow[4] = strconv.Itoa(int(taxi.DY))
			tripRow[5] = strconv.Itoa((int(taxi.DepartureTime) + int(math.Ceil(taxi.VMT/30*3600))) % 86400)
			tripRow[6] = strconv.FormatFloat(taxi.VMT, 'f', 2, 64)
			tripRow[7] = strconv.Itoa(int(taxi.NumPassengers))
			tripRow[8] = strconv.FormatFloat(taxi.PMT, 'f', 2, 64)
			oXSuper5, oYSuper5 := ataxi.GetSuperPixel(taxi.OX, taxi.OY, 3)
			dXSuper5, dYSuper5 := ataxi.GetSuperPixel(taxi.DX, taxi.DY, 3)
			oXSuper10, oYSuper10 := ataxi.GetSuperPixel(taxi.OX, taxi.OY, 4)
			dXSuper10, dYSuper10 := ataxi.GetSuperPixel(taxi.DX, taxi.DY, 4)
			tripRow[9] = strconv.Itoa(int(oXSuper5))
			tripRow[10] = strconv.Itoa(int(oYSuper5))
			tripRow[11] = strconv.Itoa(int(dXSuper5))
			tripRow[12] = strconv.Itoa(int(dYSuper5))
			tripRow[13] = strconv.Itoa(int(oXSuper10))
			tripRow[14] = strconv.Itoa(int(oYSuper10))
			tripRow[15] = strconv.Itoa(int(dXSuper10))
			tripRow[16] = strconv.Itoa(int(dYSuper10))
			tripWriter.Write(tripRow[:])
		}

		regionPMT += pmt
		regionVMT += vmt
		statePMT += pmt
		stateVMT += vmt
	}

	stateRow[0] = stateFIPS
	stateRow[1] = strconv.FormatFloat(statePMT/stateVMT, 'f', 2, 64)
    stateRow[2] = strconv.FormatFloat(statePMT, 'f', 2, 64)
    stateRow[3] = strconv.FormatFloat(stateVMT, 'f', 2, 64)
    fmt.Printf("state %s avo: %s - pmt: %s - vmt: %s\n", stateRow[0], stateRow[1], stateRow[2], stateRow[3])
	stateWriter.Write(stateRow[:])

	stateWriter.Flush()
	stateFile.Close()

	countyWriter.Flush()
	countyFile.Close()

	tripWriter.Flush()
	tripFile.Close()

	regionFile, err := os.Create("../data/region_avo.csv")
	if err != nil {
		log.Fatal(err)
	}
	regionWriter := csv.NewWriter(regionFile)
	regionColumns := []string{"Region", "AVO", "PMT", "VMT"}
	regionWriter.Write(regionColumns)
    regionAVO := regionPMT / regionVMT
	regionAVOString := strconv.FormatFloat(regionAVO, 'f', 2, 64)
    regionPMTString := strconv.FormatFloat(regionPMT, 'f', 2, 64)
    regionVMTString := strconv.FormatFloat(regionVMT, 'f', 2, 64)
	regionWriter.Write([]string{"South", regionAVOString, regionPMTString, regionVMTString})
	fmt.Printf("region avo: %s - pmt: %s - vmt: %s\n", regionAVOString, regionPMTString, regionVMTString)

	regionWriter.Flush()
	regionFile.Close()

	elapsed := time.Since(start)
	fmt.Printf("csv processing took %s\n", elapsed)
}
