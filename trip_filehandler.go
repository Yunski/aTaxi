package ataxi

import "strconv"

type Row struct {
	PersonID       int32
	OFIPS          int32
	OXCoord        int32
	OYCoord        int32
	OLat           float64
	OLon           float64
	DFIPS          int32
	DXCoord        int32
	DYCoord        int32
	DLat           float64
	DLon           float64
	ODepartureTime int32
}

func ParseLine(line []string) Row {
	OFIPS := 5
	DFIPS := 13
	PersonID := 1
	OXCoord := 8
	OYCoord := 9
	OLon := 6
	OLat := 7
	DXCoord := 16
	DYCoord := 17
	DLon := 14
	DLat := 15
	ODepartureTime := 10

	oFIPS, _ := strconv.ParseInt(line[OFIPS], 0, 32)
	dFIPS, _ := strconv.ParseInt(line[DFIPS], 0, 32)
	personID, _ := strconv.ParseInt(line[PersonID], 0, 32)
	oX, _ := strconv.ParseInt(line[OXCoord], 0, 32)
	oY, _ := strconv.ParseInt(line[OYCoord], 0, 32)
	dX, _ := strconv.ParseInt(line[DXCoord], 0, 32)
	dY, _ := strconv.ParseInt(line[DYCoord], 0, 32)
	oLat, _ := strconv.ParseFloat(line[OLat], 64)
	oLon, _ := strconv.ParseFloat(line[OLon], 64)
	dLat, _ := strconv.ParseFloat(line[DLat], 64)
	dLon, _ := strconv.ParseFloat(line[DLon], 64)
	departureTime, _ := strconv.ParseInt(line[ODepartureTime], 0, 32)

	return Row{
		PersonID:       int32(personID),
		OFIPS:          int32(oFIPS),
		OXCoord:        int32(oX),
		OYCoord:        int32(oY),
		OLat:           oLat,
		OLon:           oLon,
		DFIPS:          int32(dFIPS),
		DXCoord:        int32(dX),
		DYCoord:        int32(dY),
		DLat:           dLat,
		DLon:           dLon,
		ODepartureTime: int32(departureTime),
	}
}
