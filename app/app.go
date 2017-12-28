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

func main() {
	r := mux.NewRouter()
	r.Methods("GET").Path("/").Handler(appHandler(homeHandler))
	r.Methods("GET").Path("/api/taxis").Handler(appHandler(listTaxiHandler))
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
	fmt.Println("Listening at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// homeHandler displays the home page.
func homeHandler(w http.ResponseWriter, r *http.Request) *appError {
	w.Write([]byte("hello world"))
	return nil
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
		return appErrorf(err, "could not list taxis: %v", err)
	}
	taxiString, err := json.MarshalIndent(taxis, "", "  ")
	if err != nil {
		return appErrorf(err, "failed to return taxi json list: %v", err)
	}
	w.Write(taxiString)
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

func appErrorf(err error, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
