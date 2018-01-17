# aTaxi
Assessment of RideSharing, ‘Last-Mile” and Optimal Empty Vehicle Management of Large Regional aTaxi Operation

## About
This project establishes a computational framework to assess the ride-sharing opportunities for a given county or state.

## Requirements
[Go](https://golang.org/dl/) - After downloading, follow the installation instructions [here](https://golang.org/doc/install). \
[MySQL](https://dev.mysql.com/downloads/mysql) - Make sure to start MySQL server on your machine.

## Setup
### Config
Create a json "config.json" in the project root directory.
Project structure should look like
```
app/
config.json
```

Paste the following in "config.json", and edit the field values accordingly.
```json
{
    "username": "root",
    "password": "password",
    "database": "database_name",
    "google_maps_api_key": "your_api_key"
}
```

### Data
Create a directory "data/" in the project root directory.
Project structure should look like
```
app/
data/
config.json
```

This directory should contain your csv files (in particular [NationWide Modal Person Trip Files](http://orf467.princeton.edu/NationWideModalPersonTrips18Kyle/aTaxi/)).

To populate the MySQL database, run the following commands in terminal:
```
$ cd deploy/
$ go run db_populate.go [csv_file_name]
```

### Dependencies
Run the following commands in terminal to install Go dependencies:
```
$ go get github.com/gorilla/handlers
$ go get github.com/gorilla/mux
$ go get github.com/go-sql-driver/mysql
$ go get github.com/jinzhu/gorm
$ go get github.com/kellydunn/golang-geo
```

### Analysis
For quick generation of region analysis csv files, first run the region avo script:
```
$ cd avo/
$ go run region_avo.go path/to/modal-person-trip-files
```
This generates the `ataxi_trips.csv` file. Run the rest of the analysis scripts in the following directories:
```
cumulative/
region_totals/
supplydemand/
```
i.e.
```
$ cd supplydemand
$ go run supply_demand.go path/to/ataxi_trips.csv
```

## Server
```bash
$ cd app/
$ go run *.go
```
Server will be listening at localhost:8080.

## API
**GET** - /api/taxis \
parameters: \
**STRING** orderby = field to sort taxis (departure_time, num_passengers) \
**BOOLEAN** passengers = return passenger info along with each taxi (true, false) \
**INT** limit = number of taxis to return \
**INT** ox = X coord of origin pixel \
**INT** oy = Y coord of origin pixel \
**INT** dx_super = X coord of destination superpixel \
**INT** dy_super = Y coord of destination superpixel
