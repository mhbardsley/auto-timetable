package backend

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/mhbardsley/auto-timetable/types"
	"github.com/pelletier/go-toml/v2"
)

type event struct {
	types.Event
}

type deadline struct {
	types.Deadline
	slotsRemaining   int       `json:"-"`
	slotsAvailable   int       `json:"-"`
}

type periodic struct {
	types.Periodic
}

type inputData struct {
	Events    []event    `json:"events" toml:"events"`
	Deadlines []deadline `json:"deadlines" toml:"deadlines"`
	Periodics []periodic `json:"periodic" toml:"periodics"`
	slots     int        `json:"-"`
}

// GetInput is the function will read the JSON file into the structs
func GetInput(dirPtr *string, noOfSlots int) (data inputData) {
	tomlPaths, err := getTomls(dirPtr)
	if err != nil {
		log.Fatalf("could not find .at.toml config files: %s", err)
	}
	data, err = tomlsToInputData(tomlPaths)
	if err != nil {
		log.Fatalf("could not find any event, deadline, or periodic data: %s", err)
	}
	sortData(data)
	checkData(data)
	data.slots = noOfSlots
	return data
}

// getTomls finds all files named .at.toml in the file hierarchy, where filePtr is considered
// top of the filesystem
func getTomls(filePtr *string) ([]string, error) {
	var tomls []string
	err := filepath.WalkDir(*filePtr, func(s string, d fs.DirEntry, e error) error {
		if e != nil { return e }
		if d.Name() == ".at.toml" {
			tomls = append(tomls, s)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}
	if len(tomls) == 0 {
		return nil, fmt.Errorf("could not find any files called .at.toml")
	}
	return tomls, nil
}

// tomlsToInputData takes a list of toml files and collects them into inputData ([]events and []deadlines)
func tomlsToInputData(tomlPaths []string) (inputData, error) {
	var events []event
	var deadlines []deadline
	var periodics []periodic
	for _, tomlPath := range tomlPaths {
		dataRaw, err := os.ReadFile(tomlPath)
		if err != nil {
			log.Warnf("could not open toml file %s: %s", tomlPath, err)
			continue
		}
		var localisedInputData inputData
		err = toml.Unmarshal(dataRaw, &localisedInputData)
		if err != nil {
			log.Warnf("could not process toml file %s as valid input data: %s", tomlPath, err)
			continue
		}
		events = append(events, localisedInputData.Events...)
		deadlines = append(deadlines, localisedInputData.Deadlines...)
		periodics = append(periodics, localisedInputData.Periodics...)
	}
	if len(events) == 0 && len(deadlines) == 0 {
		return inputData{}, fmt.Errorf("could not find any events or deadlines")
	}
	return inputData{Events: events, Deadlines: deadlines, Periodics: periodics}, nil
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
		deadlines[i].DeadlineTime = roundDown(deadline.DeadlineTime)
	}
	sort.Slice(deadlines, func(p, q int) bool {
		return deadlines[p].DeadlineTime.Before(deadlines[q].DeadlineTime)
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

	// check that the event has a name
	for _, event := range events {
		if event.Name == "" {
			log.Fatalf("found an event with no name")
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
	if len(deadlines) > 0 && deadlines[0].DeadlineTime.Before(currentTime) {
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
