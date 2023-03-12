package backend

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"sort"
	"time"
)

type event struct {
	Name       string    `json:"name,omitempty"`
	Repopulate bool      `json:"repopulate,omitempty"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
}

type deadline struct {
	Name             string    `json:"name"`
	MinutesRemaining float64   `json:"minutesRemaining"`
	slotsRemaining   int       `json:"-"`
	slotsAvailable   int       `json:"-"`
	Deadline         time.Time `json:"deadline"`
}

type inputData struct {
	Events    []event    `json:"events"`
	Deadlines []deadline `json:"deadlines"`
	slots     int        `json:"-"`
}

// GetInput is the function will read the JSON file into the structs
func GetInput(filePtr *string, noOfSlots int) (data inputData) {
	dataRaw, err := os.ReadFile(*filePtr)
	if err != nil {
		log.Fatal("error opening file: ", err)
	}

	// unmarshal data into payload
	err = json.Unmarshal(dataRaw, &data)
	if err != nil {
		log.Fatal("error making sense of input file: ", err)
	}
	sortData(data)
	checkData(data)
	data.slots = noOfSlots
	return data
}

// sortData sorts events and deadlines by start date and upcoming date, respectively
// it also does rounding
func sortData(data inputData) {
	sortEvents(data.Events)
	sortDeadlines(data.Deadlines)
}

// checkData checks the validity of the data
func checkData(data inputData) {
	checkEvents(data.Events)
	checkDeadlines(data.Deadlines)
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
		deadlines[i].MinutesRemaining = math.Ceil(deadline.MinutesRemaining/25) * 25
		deadlines[i].Deadline = roundDown(deadline.Deadline)
	}
	sort.Slice(deadlines, func(p, q int) bool {
		return deadlines[p].Deadline.Before(deadlines[q].Deadline)
	})
}

// checkEvents will ensure events all start in the future, have an end date after start date
// and do not intersect
func checkEvents(events []event) {
	// if events are empty, it is trivial that they are compliant
	if len(events) == 0 {
		return
	}
	// data are sorted, so check first event
	if events[0].EndTime.Before(currentTime) {
		log.Fatalf("found an event %s that has already passed", events[0].Name)
	}

	// check that the event has a name or is to repopulate
	for _, event := range events {
		if event.Name == "" && !event.Repopulate {
			log.Fatalf("found an event that has neither a name nor is a repopulation")
		}
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
}

// checkDeadlines will ensure deadlines are in the future
func checkDeadlines(deadlines []deadline) {
	// since we assume data are sorted, just check the first deadline
	if len(deadlines) > 0 && deadlines[0].Deadline.Before(currentTime) {
		log.Fatalf("found a deadline %s that has already passed", deadlines[0].Name)
	}
}

// roundUp rounds a time up to its nearest 30-minute point
func roundUp(unrounded time.Time) (rounded time.Time) {
	rounded = unrounded.Truncate(30 * time.Minute)
	if unrounded == rounded {
		return rounded
	}
	return rounded.Add(30 * time.Minute)
}

// roundDown rounds a time down to its nearest 30-minute point
func roundDown(unrounded time.Time) time.Time {
	return unrounded.Truncate(30 * time.Minute)
}
