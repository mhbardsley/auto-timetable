package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
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

	// sort it
	sortData(data)
	return data
}

// sortData sorts events and deadlines by start date and upcoming date, respectively
func sortData(data inputData) {
	sortEvents(data.Events)
	sortDeadlines(data.Deadlines)
}

// sortEvents to sort by start time
func sortEvents(events []event) {
	sort.Slice(events, func(p, q int) bool {
		return events[p].StartTime.Before(events[q].StartTime)
	})
}

// sortDeadlines to sort by deadline
func sortDeadlines(deadlines []deadline) {
	sort.Slice(deadlines, func(p, q int) bool {
		return deadlines[p].Deadline.Before(deadlines[q].Deadline)
	})
}
