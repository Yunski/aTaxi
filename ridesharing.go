package ataxi

import (
	"log"
)

var DB RideSharingDatabase

func init() {
	var err error
	DB, err = newMySQLDB(Config)
	if err != nil {
		log.Fatal(err)
	}
}

type RideSharingDatabase interface {
	// ListTaxis returns a list of taxis, ordered by field.
	ListTaxis(orderBy string, limit int, withPassengers bool) ([]Taxi, error)

	// ListTaxisByDepartureTime returns a list of taxis, ordered by departure time.
	ListTaxisByDepartureTime(limit int, withPassengers bool) ([]Taxi, error)

	// ListTaxisByNumPassengers returns a list of taxis, ordered by number of passengers.
	ListTaxisByNumPassengers(limit int, withPassengers bool) ([]Taxi, error)

	// GetTaxi retrieves a taxi by its ID.
	GetTaxi(id uint) (*Taxi, error)

	// ListPassengers returns a list of passengers, ordered by departure time.
	ListPassengers(limit int) ([]Passenger, error)

	// GetPassenger retrieves a passenger by its ID.
	GetPassenger(id uint) (*Passenger, error)

	// GetDemandForPixels returns the supply of taxis for each pixel
	GetDemandForPixels(size int) ([]SuperPixelDemand, error)

	// GetSupplyForPixels returns the supply of taxis for each pixel
	GetSupplyForPixels(size int) ([]SuperPixelSupply, error)

	// GetNumTripsForCategory returns the number of trips for a given trip category
	GetNumTripsForCategory(category int) (int, error)

	// GetCumulativeNumTripsForCategory returns the cumulative number of trips for trip categories <= category
	GetCumulativeNumTripsForCategory(category int) (int, error)

	// Close closes the database, freeing up any available resources.
	Close()
}
