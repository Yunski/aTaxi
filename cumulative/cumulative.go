package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal(errors.New("You must provide a region totals csv."))
		os.Exit(1)
	}
	csvFile, _ := os.Open(os.Args[1])
	reader := csv.NewReader(bufio.NewReader(csvFile))
	if _, err := reader.Read(); err != nil {
		log.Fatal(err)
	}
	var totals [5]int
	line, _ := reader.Read()
	for i, field := range line {
		value, _ := strconv.ParseInt(field, 10, 64)
		if i != 0 {
			totals[i] = int(value) + totals[i-1]
		} else {
			totals[i] = int(value)
		}
	}
	fmt.Println("Creating cumulative.csv")
	file, err := os.Create("../data/cumulative.csv")
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
		row = append(row, strconv.Itoa(totals[c]))
	}
	writer.Write(row)
	fmt.Println("Finished")
}
