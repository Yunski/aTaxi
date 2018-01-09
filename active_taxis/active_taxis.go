package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal(errors.New("You must provide the generated ataxi region trips csv."))
		os.Exit(1)
	}

    csvFile, _ := os.Open(os.Args[1])
	reader := csv.NewReader(bufio.NewReader(csvFile))
	if _, err := reader.Read(); err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("../data/active_taxis.csv")
	if err != nil {
		log.Fatal(err)
	}
    defer file.Close()

    writer := csv.NewWriter(file)
	writer.Write([]string{"Min", "NumTaxis"})
    defer writer.Flush()

	start := time.Now()
	fmt.Println("Processing provided ataxi trip file ...")

    var activeTaxis [1440]int
    var counter int
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

	    departureTime, _ := strconv.ParseInt(line[2], 10, 64)
        madeEmptyTime, _ := strconv.ParseInt(line[5], 10, 64)

        start := int(departureTime) / 60
        end := int(madeEmptyTime) / 60

        activeTaxis[start]++
        if end < 1439 {
            activeTaxis[end+1]--
        }
        if start > end {
            activeTaxis[0]++
        }
        counter++
		if counter%10000 == 0 {
			fmt.Printf("\rProcessed %d records", counter)
		}
    }

    fmt.Println()

	elapsed := time.Since(start)
	fmt.Printf("ataxi trips file processing took %s\n", elapsed)

    var sum int
    for i := 0; i < len(activeTaxis); i++ {
        sum += activeTaxis[i]
        activeTaxis[i] = sum
    }

    var row [2]string
    for min, numTaxis := range activeTaxis {
        row[0] = strconv.Itoa(min)
        row[1] = strconv.Itoa(numTaxis)
        writer.Write(row[:])
    }

    fmt.Println("Successfully created active_taxis.csv")
}