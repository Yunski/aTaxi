package ataxi

import geo "github.com/kellydunn/golang-geo"

type Passenger struct {
	PersonID         int32
	OFIPS            int32
	OCoord           Coord
	OLatLon          *geo.Point
	DFIPS            int32
	DCoord           Coord
	DLatLon          *geo.Point
	DepartureTime    int32
	TripCategory     int32
	TripDistance     float64
	LatestPickUpTime int32
	DSuperPixel      Coord
}

func NewPassenger(personID int32, oFIPS int32, oCoord Coord, oLatLon *geo.Point,
	dFIPS int32, dCoord Coord, dLatLon *geo.Point, departureTime int32,
	tripCategory int32, tripDistance float64, dSuperPixel Coord) *Passenger {
	passenger := &Passenger{
		PersonID:         personID,
		OFIPS:            oFIPS,
		OCoord:           oCoord,
		OLatLon:          oLatLon,
		DFIPS:            dFIPS,
		DCoord:           dCoord,
		DLatLon:          dLatLon,
		DepartureTime:    departureTime,
		TripCategory:     tripCategory,
		TripDistance:     tripDistance,
		LatestPickUpTime: departureTime + GetMaxWaitingTime(tripDistance),
		DSuperPixel:      dSuperPixel,
	}
	return passenger
}

func (passenger *Passenger) AssignedTaxi(taxis []*Taxi) (*Taxi, bool) {
	if len(taxis) == 0 || passenger.OCoord != taxis[0].OCoord {
		return nil, true
	}
	var matchedTaxi *Taxi
	for _, taxi := range taxis {
		if taxi.DSuperPixel == passenger.DSuperPixel && taxi.HasDeparted(passenger.LatestPickUpTime) {
			matchedTaxi = taxi
			break
		}
	}
	return matchedTaxi, false
}

func (passenger *Passenger) DistanceTo(latlon *geo.Point) float64 {
	return GetTripDistance(passenger.DLatLon, latlon)
}
