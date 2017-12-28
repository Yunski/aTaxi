package ataxi

import (
	"math"
	"time"

	geo "github.com/kellydunn/golang-geo"
)

type Passenger struct {
	ID               uint `gorm:"primary_key"`
	CreatedAt        time.Time
	PersonID         int64
	OType            byte   `gorm:"index";sql:"size:1"`
	OName            string `sql:"size:255"`
	OFIPS            uint32 `gorm:"column:o_fips"`
	OX               int32
	OY               int32
	OLat             float64
	OLon             float64
	DType            byte   `gorm:"index";sql:"size:1"`
	DName            string `sql:"size:255"`
	DFIPS            uint32 `gorm:"column:d_fips"`
	DX               int32  `gorm:"index"`
	DY               int32  `gorm:"index"`
	DLat             float64
	DLon             float64
	DepartureTime    uint32 `gorm:"index"`
	TripCategory     uint32 `gorm:"index"`
	TripDistance     float64
	LatestPickUpTime uint32 `gorm:"index"`
	DXSuper          int32
	DYSuper          int32
	TaxiID           uint
}

func NewPassenger(id uint, personID int64, oType byte, oName string, oFIPS uint32,
	oX int32, oY int32, oLat float64, oLon float64, dType byte, dName string, dFIPS uint32,
	dX int32, dY int32, dLat float64, dLon float64, departureTime uint32, tripCategory uint32,
	tripDistance float64, dXSuper int32, dYSuper int32) *Passenger {
	passenger := &Passenger{
		ID:               id,
		PersonID:         personID,
		OType:            oType,
		OName:            oName,
		OFIPS:            oFIPS,
		OX:               oX,
		OY:               oY,
		OLat:             oLat,
		OLon:             oLon,
		DType:            dType,
		DName:            dName,
		DFIPS:            dFIPS,
		DX:               dX,
		DY:               dY,
		DLat:             dLat,
		DLon:             dLon,
		DepartureTime:    departureTime,
		TripCategory:     tripCategory,
		TripDistance:     tripDistance,
		LatestPickUpTime: departureTime + GetMaxWaitingTime(tripDistance),
		DXSuper:          dXSuper,
		DYSuper:          dYSuper,
	}
	return passenger
}

func NewPassengerFromRow(id uint, row Row) *Passenger {
	tripDistance := GetTripDistance(geo.NewPoint(row.OLat, row.OLon),
		geo.NewPoint(row.DLat, row.DLon))
	tripCategory := GetTripCategory(tripDistance)
	dXSuper, dYSuper := GetSuperPixel(row.DXCoord, row.DYCoord, tripCategory)

	return NewPassenger(id, row.PersonID, row.OType, row.OName, row.OFIPS,
		row.OXCoord, row.OYCoord, row.OLat, row.OLon, row.DType, row.DName,
		row.DFIPS, row.DXCoord, row.DYCoord, row.DLat, row.DLon,
		row.ODepartureTime, tripCategory, tripDistance, dXSuper, dYSuper)
}

func (passenger *Passenger) FindTaxi(taxis []*Taxi) (*Taxi, bool) {
	if len(taxis) == 0 || (passenger.OX != taxis[0].OX || passenger.OY != taxis[0].OY) {
		return nil, true
	}
	var matchedTaxi *Taxi
	for _, taxi := range taxis {
		if taxi.DXSuper == passenger.DXSuper &&
			taxi.DYSuper == passenger.DYSuper &&
			taxi.HasDeparted(passenger.LatestPickUpTime) {
			matchedTaxi = taxi
			break
		}
	}
	return matchedTaxi, false
}

func (passenger *Passenger) DistanceTo(latlon *geo.Point) float64 {
	return GetTripDistance(geo.NewPoint(passenger.DLat, passenger.DLon), latlon)
}

type Taxi struct {
	ID            uint `gorm:"primary_key"`
	CreatedAt     time.Time
	OFIPS         uint32 `gorm:"column:o_fips"`
	DFIPS         uint32 `gorm:"column:d_fips"`
	OX            int32  `gorm:"index"`
	OY            int32  `gorm:"index"`
	OLat          float64
	OLon          float64
	DX            int32 `gorm:"index"`
	DY            int32 `gorm:"index"`
	DLat          float64
	DLon          float64
	Passengers    []Passenger
	DepartureTime uint32 `gorm:"index"`
	MaxOccupancy  uint32
	NumPassengers uint32
	PMT           float64 `gorm:"column:pmt"`
	VMT           float64 `gorm:"column:vmt"`
	DXSuper       int32   `gorm:"index"`
	DYSuper       int32   `gorm:"index"`
}

func NewTaxi(id uint, passenger *Passenger, maxOccupancy uint32) *Taxi {
	taxi := &Taxi{
		ID:            id,
		OFIPS:         passenger.OFIPS,
		DFIPS:         passenger.DFIPS,
		OX:            passenger.OX,
		OY:            passenger.OY,
		OLat:          passenger.OLat,
		OLon:          passenger.OLon,
		DX:            passenger.DX,
		DY:            passenger.DY,
		DLat:          passenger.DLat,
		DLon:          passenger.DLon,
		Passengers:    []Passenger{*passenger},
		DepartureTime: passenger.LatestPickUpTime,
		MaxOccupancy:  maxOccupancy,
		NumPassengers: 1,
		DXSuper:       passenger.DXSuper,
		DYSuper:       passenger.DYSuper,
	}
	return taxi
}

func (taxi *Taxi) HasDeparted(time uint32) bool {
	return taxi.DepartureTime <= time
}

func (taxi *Taxi) AddPassenger(passenger *Passenger) {
	taxi.Passengers = append(taxi.Passengers, *passenger)
	taxi.NumPassengers++
}

func (taxi *Taxi) IsFull() bool {
	return taxi.NumPassengers == taxi.MaxOccupancy
}

func (taxi *Taxi) VehicleMilesTraveled() float64 {
	curLatLon := geo.NewPoint(taxi.OLat, taxi.OLon)
	passengers := taxi.Passengers
	var vtm float64
	naive := true
	if naive {
		for _, passenger := range passengers {
			vtm += passenger.DistanceTo(curLatLon)
			curLatLon = geo.NewPoint(passenger.DLat, passenger.DLon)
		}
	} else {
		for len(passengers) > 0 {
			nearestDist := math.Inf(1)
			nearestIdx, idx := 0, 0
			for _, passenger := range passengers {
				dist := passenger.DistanceTo(curLatLon)
				if dist < nearestDist {
					nearestDist = dist
					nearestIdx = idx
				}
				idx++
			}
			vtm += nearestDist
			curLatLon = geo.NewPoint(passengers[nearestIdx].DLat,
				passengers[nearestIdx].DLon)
			var remaining []Passenger
			for i, passenger := range passengers {
				if i != nearestIdx {
					remaining = append(remaining, passenger)
				}
			}
			passengers = remaining
		}
	}
	return vtm
}

func (taxi *Taxi) PersonMilesTraveled() float64 {
	var ptm float64
	for _, passenger := range taxi.Passengers {
		ptm += passenger.TripDistance
	}
	return ptm
}
