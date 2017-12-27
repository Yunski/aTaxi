package ataxi

type RideSharingDatabase interface {
	// ListTaxis returns a list of taxis, ordered by departure time.
	ListTaxis() ([]*Taxi, error)

	// GetTaxi retrieves a taxi by its ID.
	GetTaxi(id int32) (*Taxi, error)

	// AddTaxi saves a given taxi, assigning it a new ID.
	AddTaxi(t *Taxi) (id int32, err error)

	// DeleteTaxi removes a given taxi by its ID.
	DeleteTaxi(id int32) error

	// ListPassengers returns a list of passengers, ordered by departure time.
	ListPassengers() ([]*Passenger, error)

	// GetPassenger retrieves a passenger by its ID.
	GetPassenger(id int32) (*Taxi, error)

	// AddTaxi saves a given passenger, assigning it a new ID.
	AddPassenger(p *Passenger) (id int32, err error)

	// DeleteTaxi removes a given passenger by its ID.
	DeletePassenger(id int32) error

	// ListPassengersForTaxi returns a list of passengers for a given taxi, ordered by departure time.
	ListPassengersForTaxi(t *Taxi) ([]*Passenger, error)

	// Close closes the database, freeing up any available resources.
	// TODO(cbro): Close() should return an error.
	Close()
}
