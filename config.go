package ataxi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type AppConfig struct {
	Username         string
	Password         string
	Database         string
	GoogleMapsAPIKey string `json:"google_maps_api_key"`
}

var Config AppConfig

func init() {
	raw, err := ioutil.ReadFile("../config.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	json.Unmarshal(raw, &Config)
}
