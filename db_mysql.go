package ataxi

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type mysqlDB struct {
	conn *gorm.DB
}

var _ RideSharingDatabase = &mysqlDB{}

func newMySQLDB(config AppConfig) (RideSharingDatabase, error) {
	conn, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8&parseTime=True&loc=Local", config.Username, config.Password, config.Database))
	if err != nil {
		return nil, fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	db := &mysqlDB{
		conn: conn,
	}
	return db, nil
}

// ListTaxis returns a list of taxis, ordered by field.
func (db *mysqlDB) ListTaxis(orderBy string, limit int, withPassengers bool) ([]Taxi, error) {
	var taxis []Taxi
	var err error
	if orderBy == "departure_time" {
		taxis, err = db.ListTaxisByDepartureTime(limit, withPassengers)
	} else {
		taxis, err = db.ListTaxisByNumPassengers(limit, withPassengers)
	}
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return taxis, err
}

// ListTaxis returns a list of taxis, ordered by departure time.
func (db *mysqlDB) ListTaxisByDepartureTime(limit int, withPassengers bool) ([]Taxi, error) {
	var taxis []Taxi
	if withPassengers {
		db.conn.Limit(limit).Preload("Passengers").Order("departure_time asc").Find(&taxis)
	} else {
		db.conn.Limit(limit).Order("departure_time asc").Find(&taxis)
	}
	if len(taxis) != limit {
		return nil, errors.New("mysql: could not retrieve taxis")
	}
	return taxis, nil
}

// ListTaxis returns a list of taxis, ordered by number of passengers.
func (db *mysqlDB) ListTaxisByNumPassengers(limit int, withPassengers bool) ([]Taxi, error) {
	var taxis []Taxi
	if withPassengers {
		db.conn.Limit(limit).Preload("Passengers").Order("num_passengers desc").Find(&taxis)
	} else {
		db.conn.Limit(limit).Order("num_passengers desc").Find(&taxis)
	}
	if len(taxis) != limit {
		return nil, errors.New("mysql: could not retrieve taxis")
	}
	return taxis, nil
}

// GetTaxi retrieves a taxi by its ID.
func (db *mysqlDB) GetTaxi(id uint) (*Taxi, error) {
	var taxi Taxi
	db.conn.Preload("Passengers").First(&taxi, id)
	if taxi.ID == 0 {
		return nil, fmt.Errorf("mysql: could not find taxi with id %d", id)
	}
	return &taxi, nil
}

// ListPassengers returns a list of passengers, ordered by departure time.
func (db *mysqlDB) ListPassengers(limit int) ([]Passenger, error) {
	var passengers []Passenger
	db.conn.Limit(limit).Order("departure_time asc").Find(&passengers)
	if len(passengers) != limit {
		return nil, errors.New("mysql: could not retrieve passengers")
	}
	return passengers, nil
}

// GetPassenger retrieves a passenger by its ID.
func (db *mysqlDB) GetPassenger(id uint) (*Passenger, error) {
	var passenger Passenger
	db.conn.First(&passenger, id)
	if passenger.ID == 0 {
		return nil, fmt.Errorf("mysql: could not find passenger with id %d", id)
	}
	return &passenger, nil
}

// GetDemandForPixels returns the supply of taxis for each pixel
func (db *mysqlDB) GetDemandForPixels(size int) ([]SuperPixelDemand, error) {
	if size != 1 && size != 5 && size != 10 {
		return nil, fmt.Errorf("Superpixel of dimension %dx%d is not supported", size)
	}
	var results []SuperPixelDemand
	db.conn.Raw("select count(*) as c, ox, oy from taxis group by ox, oy").Scan(&results)
	return results, nil
}

// GetSupplyForPixels returns the supply of taxis for each pixel
func (db *mysqlDB) GetSupplyForPixels(size int) ([]SuperPixelSupply, error) {
	if size != 1 && size != 5 && size != 10 {
		return nil, fmt.Errorf("Superpixel of dimension %dx%d is not supported", size)
	}
	var results []SuperPixelSupply
	db.conn.Raw("select count(*) as c, dx_super, dy_super from taxis group by dx_super, dy_super").Scan(&results)
	return results, nil
}

func (db *mysqlDB) GetNumTripsForCategory(category int) (int, error) {
	var numTrips int
	db.conn.Model(&Passenger{}).Where("trip_category = ?", category).Count(&numTrips)
	return numTrips, nil
}

func (db *mysqlDB) GetCumulativeNumTripsForCategory(category int) (int, error) {
	var numTrips int
	db.conn.Model(&Passenger{}).Where("trip_category <= ?", category).Count(&numTrips)
	return numTrips, nil
}

// Close closes the database, freeing up any available resources.
func (db *mysqlDB) Close() {
	db.conn.Close()
}
