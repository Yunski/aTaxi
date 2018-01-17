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
    const mileBuckets = 401
    var tripLengthCumulative [mileBuckets]int
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
            lengthIdx := int(math.Floor(tripDistance))
            if lengthIdx >= mileBuckets-1 {
                lengthIdx = mileBuckets-1
            }
            tripLengthCumulative[lengthIdx]++
			tripCategory := ataxi.GetTripCategory(tripDistance)
			regionTotals[int(tripCategory)]++
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("csv processing took %s\n", elapsed)
	fmt.Println("Creating region_totals.csv")
	regionFile, err := os.Create("../data/region_totals.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer regionFile.Close()

	regionWriter := csv.NewWriter(regionFile)
	defer regionWriter.Flush()

	columns := []string{"WalkTrips", "ShortTrips", "NormalTrips", "LongTrips", "ReallyLongTrips"}
	regionWriter.Write(columns)

	var row []string
	for c := 0; c < 5; c++ {
		row = append(row, strconv.Itoa(regionTotals[c]))
	}
	regionWriter.Write(row)
	fmt.Println("Finished")

    for i := 0; i < len(tripLengthCumulative); i++ {
        if i != 0 {
            tripLengthCumulative[i] += tripLengthCumulative[i-1]
        }
    }
    fmt.Println("Creating trip_length_cumulative.csv")
    tripLengthFile, err := os.Create("../data/trip_length_cumulative.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer tripLengthFile.Close()

    tripLengthWriter := csv.NewWriter(tripLengthFile)
    defer tripLengthWriter.Flush()

    columns = []string{"mile", "cumulative"}
    tripLengthWriter.Write(columns)

    for mile, count := range tripLengthCumulative {
        row[0] = strconv.Itoa(mile)
        row[1] = strconv.Itoa(count)
        tripLengthWriter.Write(row[:2])
    }
    fmt.Println("Finished")
}
