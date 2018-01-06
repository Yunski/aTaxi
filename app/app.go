package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/webapps/ataxi"
)

var (
	mapTmpl = parseTemplate("map.html")
)

func main() {
	r := mux.NewRouter()
	r.Methods("GET").Path("/").Handler(appHandler(homeHandler))
	r.Methods("GET").Path("/api/taxis").Handler(appHandler(listTaxiHandler))
	r.Methods("GET").Path("/api/taxis/num_trips").Handler(appHandler(numTripsForCategoryHandler))
	r.Methods("GET").Path("/api/taxis/supply_demand").Handler(appHandler(supplyAndDemandHandler))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
	fmt.Println("Listening at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// homeHandler displays the home page.
func homeHandler(w http.ResponseWriter, r *http.Request) *appError {
	return mapTmpl.Execute(w, r, nil)
}

// listTaxiHandler returns a json list of taxis sorted by request field.
func listTaxiHandler(w http.ResponseWriter, r *http.Request) *appError {
	params := r.URL.Query()
	orderBy := "departure_time"
	if orderParam, ok := params["orderby"]; ok {
		orderBy = orderParam[0]
		if orderBy != "departure_time" {
			orderBy = "num_passengers"
		}
	}
	limit := 100
	if limitParam, ok := params["limit"]; ok {
		limit64, err := strconv.ParseInt(limitParam[0], 10, 32)
		if err == nil {
			if limit64 < 100 {
				limit = int(limit64)
			}
		}
	}
	var withPassengers bool
	if passengersParam, ok := params["passengers"]; ok {
		withPassengers, _ = strconv.ParseBool(passengersParam[0])
	}
	taxis, err := ataxi.DB.ListTaxis(orderBy, limit, withPassengers)
	if err != nil {
		return appErrorf(err, 500, "could not list taxis: %v", err)
	}
	jsonOutput, err := json.MarshalIndent(taxis, "", "  ")
	if err != nil {
		return appErrorf(err, 500, "failed to return taxi json data: %v", err)
	}
	w.Write(jsonOutput)
	return nil
}

func supplyAndDemandHandler(w http.ResponseWriter, r *http.Request) *appError {
	demandResults, err := ataxi.DB.GetDemandForPixels(1)
	if err != nil {
		return appErrorf(err, 500, "could not list demand for pixels: %v", err)
	}
	supplyResults, err := ataxi.DB.GetSupplyForPixels(1)
	if err != nil {
		return appErrorf(err, 500, "could not list supply for pixels: %v", err)
	}
	netTaxis := make(map[int32]*ataxi.SuperPixelDemand)
	for _, result := range demandResults {
		spd := &ataxi.SuperPixelDemand{
			Count: -result.Count,
			X:     result.X,
			Y:     result.Y,
		}
		netTaxis[ataxi.HashCode(result.X, result.Y)] = spd
	}
	for _, result := range supplyResults {
		spd := &ataxi.SuperPixelDemand{
			Count: result.Count,
			X:     result.X,
			Y:     result.Y,
		}
		netTaxis[ataxi.HashCode(result.X, result.Y)] = spd
	}
	var supplyDemand []ataxi.SuperPixelDemand
	for _, sd := range netTaxis {
		supplyDemand = append(supplyDemand, *sd)
	}
	jsonOutput, err := json.MarshalIndent(supplyDemand, "", "  ")
	if err != nil {
		return appErrorf(err, 500, "failed to return taxi json data: %v", err)
	}
	w.Write(jsonOutput)
	return nil
}

func numTripsForCategoryHandler(w http.ResponseWriter, r *http.Request) *appError {
	params := r.URL.Query()
	var category int
	if categoryParam, ok := params["category"]; ok {
		category64, err := strconv.ParseInt(categoryParam[0], 10, 32)
		if err != nil {
			return appErrorf(err, 422, "category param does not contain an int: \"%s\"", categoryParam[0])
		}
		category = int(category64)
	}
	if category < 1 || category > 4 {
		return appErrorf(nil, 400, "invalid trip category int")
	}
	var cumulative bool
	var err error
	if cumulativeParam, ok := params["cumulative"]; ok {
		cumulative, err = strconv.ParseBool(cumulativeParam[0])
		if err != nil {
			return appErrorf(err, 422, "cumulative is a boolean param: \"%s\" provided", cumulativeParam[0])
		}
	}
	var numTrips int
	if cumulative {
		numTrips, err = ataxi.DB.GetCumulativeNumTripsForCategory(category)
	} else {
		numTrips, err = ataxi.DB.GetNumTripsForCategory(category)
	}
	if err != nil {
		return appErrorf(err, 500, "failed to return number of trips: %v", err)
	}
	res := make(map[string]int)
	res["trip_category"] = category
	res["num_trips"] = numTrips
	jsonOutput, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return appErrorf(err, 500, "failed to return num trips json data: %v", err)
	}
	w.Write(jsonOutput)
	return nil
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

func appErrorf(err error, code int, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Code:    code,
		Message: fmt.Sprintf(format, v...),
	}
}
