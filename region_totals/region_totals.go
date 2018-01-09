package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	geo "github.com/kellydunn/golang-geo"
	"github.com/webapps/ataxi"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal(errors.New("You must provide a data directory containing the ataxi mode trip files."))
		os.Exit(1)
	}

	files, err := filepath.Glob(filepath.Join(os.Args[1], "*.csv"))
	if err != nil {
		log.Fatal(err)
	}

	regionTotals := make(map[int]int)
	start := time.Now()
	for _, file := range files {
		_, filename := filepath.Split(file)
		fmt.Printf("Processing %s\n", filename)
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
			tripDistance := ataxi.GetTripDistance(geo.NewPoint(row.OLat, row.OLon),
				geo.NewPoint(row.DLat, row.DLon))
			tripCategory := ataxi.GetTripCategory(tripDistance)
			regionTotals[int(tripCategory)]++
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("csv processing took %s\n", elapsed)
	fmt.Println("Creating region_totals.csv")
	file, err := os.Create("../data/region_totals.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	columns := []string{"WalkTrips", "ShortTrips", "NormalTrips", "LongTrips", "ReallyLongTrips"}
	writer.Write(columns)

	var row []string
	for c := 0; c < 5; c++ {
		row = append(row, strconv.Itoa(regionTotals[c]))
	}
	writer.Write(row)
	fmt.Println("Finished")
}
