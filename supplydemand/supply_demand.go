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

	"github.com/webapps/ataxi"
)

func process(X int32, Y int32, netTaxis map[int32]*ataxi.SuperPixelDemand, isSupply bool) {
	hash := ataxi.HashCode(X, Y)
	if spd, ok := netTaxis[hash]; ok {
		if isSupply {
			spd.Count++
		} else {
			spd.Count--
		}
	} else {
		spd := &ataxi.SuperPixelDemand{
			X: X,
			Y: Y,
		}
		if isSupply {
			spd.Count = 1
		} else {
			spd.Count = -1
		}
		netTaxis[hash] = spd
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(errors.New("You must provide the generated ataxi region trips csv."))
		os.Exit(1)
	}

	columns := []string{"X", "Y", "net"}
	pixelsFile, err := os.Create("../data/supplydemand_1x1.csv")
	if err != nil {
		log.Fatal(err)
	}

	pixelsWriter := csv.NewWriter(pixelsFile)
	pixelsWriter.Write(columns)

	super5x5File, err := os.Create("../data/supplydemand_5x5.csv")
	if err != nil {
		log.Fatal(err)
	}
	super5x5Writer := csv.NewWriter(super5x5File)
	super5x5Writer.Write(columns)

	super10x10File, err := os.Create("../data/supplydemand_10x10.csv")
	if err != nil {
		log.Fatal(err)
	}
	super10x10Writer := csv.NewWriter(super10x10File)
	super10x10Writer.Write(columns)

	netTaxis1x1 := make(map[int32]*ataxi.SuperPixelDemand)
	netTaxis5x5 := make(map[int32]*ataxi.SuperPixelDemand)
	netTaxis10x10 := make(map[int32]*ataxi.SuperPixelDemand)

	start := time.Now()
	fmt.Println("Processing provided ataxi trip file ...")

	csvFile, _ := os.Open(os.Args[1])
	reader := csv.NewReader(bufio.NewReader(csvFile))
	if _, err := reader.Read(); err != nil {
		log.Fatal(err)
	}
	var counter int
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		OX, _ := strconv.ParseInt(line[0], 10, 64)
		OY, _ := strconv.ParseInt(line[1], 10, 64)
		DX, _ := strconv.ParseInt(line[3], 10, 64)
		DY, _ := strconv.ParseInt(line[4], 10, 64)
		OXSuper5, _ := strconv.ParseInt(line[9], 10, 64)
		OYSuper5, _ := strconv.ParseInt(line[10], 10, 64)
		DXSuper5, _ := strconv.ParseInt(line[11], 10, 64)
		DYSuper5, _ := strconv.ParseInt(line[12], 10, 64)
		OXSuper10, _ := strconv.ParseInt(line[13], 10, 64)
		OYSuper10, _ := strconv.ParseInt(line[14], 10, 64)
		DXSuper10, _ := strconv.ParseInt(line[15], 10, 64)
		DYSuper10, _ := strconv.ParseInt(line[16], 10, 64)

		process(int32(OX), int32(OY), netTaxis1x1, false)
		process(int32(DX), int32(DY), netTaxis1x1, true)
		process(int32(OXSuper5), int32(OYSuper5), netTaxis5x5, false)
		process(int32(DXSuper5), int32(DYSuper5), netTaxis5x5, true)
		process(int32(OXSuper10), int32(OYSuper10), netTaxis10x10, false)
		process(int32(DXSuper10), int32(DYSuper10), netTaxis10x10, true)
		counter++
		if counter%10000 == 0 {
			fmt.Printf("\rProcessed %d records", counter)
		}
	}

    fmt.Println()

	var row [3]string
	for _, spd := range netTaxis1x1 {
		row[0] = strconv.Itoa(int(spd.X))
		row[1] = strconv.Itoa(int(spd.Y))
		row[2] = strconv.Itoa(spd.Count)
		pixelsWriter.Write(row[:])
	}
	for _, spd := range netTaxis5x5 {
		row[0] = strconv.Itoa(int(spd.X))
		row[1] = strconv.Itoa(int(spd.Y))
		row[2] = strconv.Itoa(spd.Count)
		super5x5Writer.Write(row[:])
	}
	for _, spd := range netTaxis10x10 {
		row[0] = strconv.Itoa(int(spd.X))
		row[1] = strconv.Itoa(int(spd.Y))
		row[2] = strconv.Itoa(spd.Count)
		super10x10Writer.Write(row[:])
	}

	elapsed := time.Since(start)
	fmt.Printf("ataxi trips file processing took %s\n", elapsed)

	pixelsWriter.Flush()
	pixelsFile.Close()

	super5x5Writer.Flush()
	super5x5File.Close()

	super10x10Writer.Flush()
	super10x10File.Close()
}
