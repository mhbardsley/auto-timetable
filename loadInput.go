package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

type event struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type deadline struct {
	Name             string    `json:"name"`
	MinutesRemaining float64   `json:"minutesRemaining"`
	Deadline         time.Time `json:"deadline"`
}

type inputData struct {
	Events    []event    `json:"events"`
	Deadlines []deadline `json:"deadlines"`
}

// this function will read the JSON file into the structs
func getInput(filePtr *string) (data inputData) {
	dataRaw, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		log.Fatal("error opening file: ", err)
	}

	// unmarshall data into payload
	err = json.Unmarshal(dataRaw, &data)
	if err != nil {
		log.Fatal("error making sense of input file: ", err)
	}
	return data
}
