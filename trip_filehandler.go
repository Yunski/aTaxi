package ataxi

import "strconv"

const PersonID = 1
const OType = 3
const OName = 4
const OFIPS = 5
const OLon = 6
const OLat = 7
const OXCoord = 8
const OYCoord = 9
const ODepartureTime = 10
const DType = 11
const DName = 12
const DFIPS = 13
const DLon = 14
const DLat = 15
const DXCoord = 16
const DYCoord = 17

type Row struct {
	PersonID       int64
	OType          byte
	OName          string
	OFIPS          uint32
	OLon           float64
	OLat           float64
	OXCoord        int32
	OYCoord        int32
	ODepartureTime uint32
	DType          byte
	DName          string
	DFIPS          uint32
	DLon           float64
	DLat           float64
	DXCoord        int32
	DYCoord        int32
}

func ParseLine(line []string) Row {
	personID, _ := strconv.ParseInt(line[PersonID], 10, 64)
	oType := byte(line[OType][0])
	oName := line[OName]
	oFIPS, _ := strconv.ParseUint(line[OFIPS], 10, 32)
	oLon, _ := strconv.ParseFloat(line[OLon], 64)
	oLat, _ := strconv.ParseFloat(line[OLat], 64)
	oX, _ := strconv.ParseInt(line[OXCoord], 10, 32)
	oY, _ := strconv.ParseInt(line[OYCoord], 10, 32)
	oDepartureTime, _ := strconv.ParseUint(line[ODepartureTime], 10, 32)
	dType := byte(line[DType][0])
	dName := line[DName]
	dFIPS, _ := strconv.ParseUint(line[DFIPS], 10, 32)
	dLon, _ := strconv.ParseFloat(line[DLon], 64)
	dLat, _ := strconv.ParseFloat(line[DLat], 64)
	dX, _ := strconv.ParseInt(line[DXCoord], 10, 32)
	dY, _ := strconv.ParseInt(line[DYCoord], 10, 32)

	return Row{
		PersonID:       personID,
		OType:          oType,
		OName:          oName,
		OFIPS:          uint32(oFIPS),
		OLon:           oLon,
		OLat:           oLat,
		OXCoord:        int32(oX),
		OYCoord:        int32(oY),
		ODepartureTime: uint32(oDepartureTime),
		DType:          dType,
		DName:          dName,
		DFIPS:          uint32(dFIPS),
		DLon:           dLon,
		DLat:           dLat,
		DXCoord:        int32(dX),
		DYCoord:        int32(dY),
	}
}
