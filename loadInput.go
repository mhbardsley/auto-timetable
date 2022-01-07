package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"sort"
	"time"
)

type event struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type reducedEvent struct {
	startTime time.Time
	endTime   time.Time
}

type deadline struct {
	Name             string    `json:"name"`
	MinutesRemaining float64   `json:"minutesRemaining"`
	slotsRemaining   int       `json:"-"`
	slotsAvailable   int       `json:"-"`
	Deadline         time.Time `json:"deadline"`
}

type reducedDeadline struct {
	minutesRemaining float64
	deadline         time.Time
}

type inputData struct {
	Events       []event                    `json:"events"`
	Deadlines    []deadline                 `json:"deadlines"`
	eventsMap    map[string]reducedEvent    `json:"-"`
	deadlinesMap map[string]reducedDeadline `json:"-"`
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

	sortData(data)
	// populate the maps
	populateMaps(&data)
	checkData(data)
	return data
}

// populateMaps populates eventsMap and deadlinesMap with the relevant data from slices already populated
func populateMaps(data *inputData) {
	data.eventsMap = make(map[string]reducedEvent)
	data.deadlinesMap = make(map[string]reducedDeadline)
	for _, event := range data.Events {
		data.eventsMap[event.Name] = reducedEvent{event.StartTime, event.EndTime}
	}
	for _, deadline := range data.Deadlines {
		data.deadlinesMap[deadline.Name] = reducedDeadline{deadline.MinutesRemaining, deadline.Deadline}
	}
}

// sortData sorts events and deadlines by start date and upcoming date, respectively
// it also does rounding
func sortData(data inputData) {
	sortEvents(data.Events)
	sortDeadlines(data.Deadlines)
}

// checkData checks the validity of the data
func checkData(data inputData) {
	checkEvents(data.Events, data.eventsMap)
	checkDeadlines(data.Deadlines, data.deadlinesMap)
}

// sortEvents to sort by start time
func sortEvents(events []event) {
	for i, event := range events {
		events[i].StartTime = roundDown(event.StartTime)
		events[i].EndTime = roundUp(event.EndTime)
	}
	sort.Slice(events, func(p, q int) bool {
		return events[p].StartTime.Before(events[q].StartTime)
	})
}

// sortDeadlines to sort by deadline
func sortDeadlines(deadlines []deadline) {
	for i, deadline := range deadlines {
		deadlines[i].MinutesRemaining = math.Ceil(deadline.MinutesRemaining/30) * 30
		deadlines[i].Deadline = roundDown(deadline.Deadline)
	}
	sort.Slice(deadlines, func(p, q int) bool {
		return deadlines[p].Deadline.Before(deadlines[q].Deadline)
	})
}

// checkEvents will ensure events all start in the future, have an end date after start date
// and do not intersect
func checkEvents(events []event, eventsMap map[string]reducedEvent) {
	// if events are empty, it is trivial that they are compliant
	if len(events) == 0 {
		return
	}
	// data are sorted, so check first event
	if events[0].EndTime.Before(currentTime) {
		log.Fatalf("found an event %s that has already passed", events[0].Name)
	}

	// check every event's start time is before the end time
	for _, event := range events {
		if event.EndTime.Before(event.StartTime) {
			log.Fatalf("found an event %s with end time before start time", event.Name)
		}
	}

	// now check for each event that it's successor does not start before it
	for i := range events {
		if i >= len(events)-1 {
			break
		}
		if events[i+1].StartTime.Before(events[i].EndTime) {
			log.Fatalf("found an event %s with start time before event %s ends", events[i+1].Name, events[i].Name)
		}
	}

	// finally check that the lengths of the maps and slices are the same
	if len(events) != len(eventsMap) {
		log.Fatal("there are events with the same name")
	}
}

// checkDeadlines will ensure deadlines are in the future
func checkDeadlines(deadlines []deadline, deadlinesMap map[string]reducedDeadline) {
	// since we assume data are sorted, just check the first deadline
	if len(deadlines) > 0 && deadlines[0].Deadline.Before(currentTime) {
		log.Fatalf("found a deadline %s that has already passed", deadlines[0].Name)
	}
	// finally check that the lengths of the maps and slices are the same
	if len(deadlines) != len(deadlinesMap) {
		log.Fatal("there are deadlines with the same name")
	}
}
