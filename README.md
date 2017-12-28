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
    "database": "database_name"
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

This directory should contain your csv files (in particular [NationWide Trip Files](http://orf467.princeton.edu/NationWideTrips'18Kyle/)).

To populate the MySQL database, run the following commands in terminal:
```bash
cd deploy/
go run db_populate.go [csv_file_name]
```

### Dependencies
Run the following commands in terminal to install Go dependencies:
```bash
go get github.com/gorilla/handlers
go get github.com/gorilla/mux
go get github.com/go-sql-driver/mysql
go get github.com/jinzhu/gorm
```

## Start server
```bash
cd app/
go run app.go
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
