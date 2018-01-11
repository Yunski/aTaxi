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

func main() {
	if len(os.Args) != 2 {
		log.Fatal(errors.New("You must provide the generated ataxi region trips csv."))
		os.Exit(1)
	}

    timeCategoryColumns := []string{"Overning", "MorningPeak", "MorningLull", "EarlyAfternoon", "EveningRush", "Evening"}
	timeCategoryFile, err := os.Create("../data/trip_distribution_time_categories.csv")
	if err != nil {
		log.Fatal(err)
	}

	timeCategoryWriter := csv.NewWriter(timeCategoryFile)
	timeCategoryWriter.Write(timeCategoryColumns)

    timeCategoryAVOFile, err := os.Create("../data/avo_time_categories.csv")
    if err != nil {
		log.Fatal(err)
	}

    timeCategoryAVOWriter := csv.NewWriter(timeCategoryAVOFile)
    timeCategoryAVOWriter.Write(timeCategoryColumns)

    hoursFile, err := os.Create("../data/trip_distribution_hours.csv")
    if err != nil {
		log.Fatal(err)
	}

    var hourColumns []string
    for i := 0; i < 24; i++ {
        hourColumns = append(hourColumns, strconv.Itoa(i))
    }

    hoursWriter := csv.NewWriter(hoursFile)
    hoursWriter.Write(hourColumns)

    hoursAVOFile, err := os.Create("../data/avo_hours.csv")
    if err != nil {
		log.Fatal(err)
	}

    hoursAVOWriter := csv.NewWriter(hoursAVOFile)
    hoursAVOWriter.Write(hourColumns)

	start := time.Now()
	fmt.Println("Processing provided ataxi trip file ...")

	csvFile, _ := os.Open(os.Args[1])
	reader := csv.NewReader(bufio.NewReader(csvFile))
	if _, err := reader.Read(); err != nil {
		log.Fatal(err)
	}

    var tripDistributionCategories [6]int
    var tripDistributionHours [24]int
    var vmtCategories [6]float64
    var pmtCategories [6]float64
    var vmtHours [24]float64
    var pmtHours [24]float64

	var counter int
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		departureTime, _ := strconv.ParseInt(line[2], 10, 64)
        vmt, _ := strconv.ParseFloat(line[6], 64)
        pmt, _ := strconv.ParseFloat(line[8], 64)
        numPassengers, _ := strconv.ParseInt(line[7], 10, 64)

        category := ataxi.GetTimeCategory(int(departureTime))
        hour := ataxi.GetHour(int(departureTime))

        tripDistributionCategories[category] += int(numPassengers)
        tripDistributionHours[hour] += int(numPassengers)
        vmtCategories[category] += vmt
        pmtCategories[category] += pmt
        vmtHours[hour] += vmt
        pmtHours[hour] += pmt

		counter++
		if counter%10000 == 0 {
			fmt.Printf("\rProcessed %d records", counter)
		}
	}

    fmt.Println()
    elapsed := time.Since(start)
	fmt.Printf("ataxi trips file processing took %s\n", elapsed)

    var avoCategories [6]float64
    for i := 0; i < 6; i++ {
        avoCategories[i] = pmtCategories[i] / vmtCategories[i]
    }
    var avoHours [24]float64
    for i := 0; i < 24; i++ {
        avoHours[i] = pmtHours[i] / vmtHours[i]
    }

    var tripDistributionCategoriesRow [6]string
    var avoCategoriesRow [6]string
    var tripDistributionHoursRow [24]string
    var avoHoursRow [24]string

    for i := 0; i < 6; i++ {
        tripDistributionCategoriesRow[i] = strconv.Itoa(tripDistributionCategories[i])
        avoCategoriesRow[i] = strconv.FormatFloat(avoCategories[i], 'f', 2, 64)
    }

    for i := 0; i < 24; i++ {
        tripDistributionHoursRow[i] = strconv.Itoa(tripDistributionHours[i])
        avoHoursRow[i] = strconv.FormatFloat(avoHours[i], 'f', 2, 64)
    }

    timeCategoryWriter.Write(tripDistributionCategoriesRow[:])
    timeCategoryAVOWriter.Write(avoCategoriesRow[:])

    hoursWriter.Write(tripDistributionHoursRow[:])
    hoursAVOWriter.Write(avoHoursRow[:])

    timeCategoryWriter.Flush()
    timeCategoryFile.Close()
	timeCategoryAVOWriter.Flush()
    timeCategoryAVOFile.Close()
    hoursWriter.Flush()
    hoursFile.Close()
    hoursAVOWriter.Flush()
    hoursAVOFile.Close()
}

