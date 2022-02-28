package main

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

type timetableElement struct {
	event    *event
	deadline *deadline
}

// function to generate the timetable
func generateTimetable(data inputData, threshold float64) {
	timetable := getEmptyTimetable(data.Deadlines, data.Events, data.slots)

	fillWithEvents(timetable, data.Events)

	if !repopulationExceeded(timetable, threshold) {
		log.Fatalf("You need to add more repopulation events/make more existing events repopulation events")
	}

	fillDeadlines(timetable, data.Deadlines)

	// if a timetabling is not possible, stop
	if time, possible := possibleTimetabling(data.Deadlines); !possible {
		log.Fatalf("There's too little time to do everything before %s! Please reduce the number of events or deadlines or extend them", time.Format("Jan 2 15:04"))
	}

	// otherwise, we loop in a random to probabilistic assignment
	fillTimetable(timetable, data.Deadlines)

	timetable = extendTimetable(timetable, data.slots)
	fmt.Printf("%s", printTimetable(timetable, data.slots))

}

// generate a slice of timetable elements
func getEmptyTimetable(deadlines []deadline, events []event, noOfSlots int) (timetable []timetableElement) {
	var numberOfSpaces int
	deadlinesEnd := len(deadlines)
	eventsEnd := len(events)
	if deadlinesEnd != 0 {
		deadlinesEnd = segmentsBetween(currentTime, deadlines[deadlinesEnd-1].Deadline)
	}
	if eventsEnd != 0 {
		eventsEnd = segmentsBetween(currentTime, events[eventsEnd-1].EndTime)
	}
	numberOfSpaces = int(math.Max(float64(deadlinesEnd), float64(eventsEnd)))
	timetable = make([]timetableElement, numberOfSpaces)
	return timetable
}

// return the number of segments between time1 (rounded up) and time2 (rounded down)
func segmentsBetween(time1 time.Time, time2 time.Time) int {
	durationBetween := time2.Sub(time1)
	return (int)(durationBetween.Minutes() / 30)
}

// fill the timetable with the events now they are assumed to be correct
func fillWithEvents(timetable []timetableElement, events []event) {
	var startIndex int
	var endIndex int
	var selectedElements []timetableElement
	for i, event := range events {
		// if the event is sufficiently late, break
		if segmentsBetween(currentTime, event.StartTime) >= len(timetable) {
			break
		}
		startIndex = int(math.Max(0, float64(segmentsBetween(currentTime, event.StartTime))))
		endIndex = int(math.Min(float64(segmentsBetween(currentTime, event.EndTime)), float64(len(timetable))))
		selectedElements = timetable[startIndex:endIndex]
		for j := range selectedElements {
			selectedElements[j].event = &(events[i])
		}
	}
}

// check if the repopulation threshold is being exceeded at any point
func repopulationExceeded(timetable []timetableElement, threshold float64) bool {
	slotsTotal := 0
	slotsRepopulation := 0
	for _, slot := range timetable {
		slotsTotal++
		if slot.event != nil && slot.event.Repopulate {
			slotsRepopulation++
		}
		if float64(slotsRepopulation)/float64(slotsTotal) > threshold {
			return true
		}
	}
	return false
}

// function to fill deadlines with how many remain and are available
func fillDeadlines(timetable []timetableElement, deadlines []deadline) {
	var startIndex int
	var endIndex int
	var currentSlots int
	startIndex = 0
	runningTotal := 0
	for i, deadline := range deadlines {
		deadlines[i].slotsRemaining = int(deadline.MinutesRemaining / 25)
		endIndex = segmentsBetween(currentTime, deadline.Deadline)
		currentSlots = freeSlotsBetween(timetable[startIndex:endIndex])
		runningTotal += currentSlots
		deadlines[i].slotsAvailable = runningTotal
		startIndex = endIndex
	}
}

// check that a timetabling is possible
func possibleTimetabling(deadlines []deadline) (noFit time.Time, possible bool) {
	runningTotal := 0
	for _, deadline := range deadlines {
		runningTotal += deadline.slotsRemaining
		if runningTotal > deadline.slotsAvailable {
			return deadline.Deadline, false
		}
	}
	return noFit, true
}

// calculate the number of free slots in the timetable slice
func freeSlotsBetween(timetablePart []timetableElement) int {
	count := 0
	for _, slot := range timetablePart {
		if slot.event == nil {
			count++
		}
	}
	return count
}

// function to extend the timetable, if need-be, to the number of slots given
func extendTimetable(timetable []timetableElement, noOfSlots int) []timetableElement {
	if noOfSlots <= len(timetable) {
		return timetable
	}
	// otherwise, make the slice
	extraPart := make([]timetableElement, noOfSlots-len(timetable))
	return append(timetable, extraPart...)
}

// function to print the timetable as-is
func printTimetable(timetable []timetableElement, noOfSlots int) string {
	builder := strings.Builder{}
	for i := 0; i < noOfSlots; i++ {
		switch {
		case timetable[i].event != nil && timetable[i].event.Repopulate && timetable[i].event.Name != "":
			builder.WriteString(fmt.Sprintf("%s-%s: ", (currentTime.Add(time.Duration(i*30) * time.Minute)).Format("Jan 2 15:04"), (currentTime.Add(time.Duration((i+1)*30) * time.Minute)).Format("Jan 2 15:04")))
			builder.WriteString(fmt.Sprintf("{REPOPULATE DEADLINES, [EVENT] %s}", timetable[i].event.Name))
		case timetable[i].event != nil && timetable[i].event.Repopulate:
			builder.WriteString(fmt.Sprintf("%s-%s: ", (currentTime.Add(time.Duration(i*30) * time.Minute)).Format("Jan 2 15:04"), (currentTime.Add(time.Duration((i+1)*30) * time.Minute)).Format("Jan 2 15:04")))
			builder.WriteString("REPOPULATE DEADLINES")
		case timetable[i].event != nil:
			builder.WriteString(fmt.Sprintf("%s-%s: ", (currentTime.Add(time.Duration(i*30) * time.Minute)).Format("Jan 2 15:04"), (currentTime.Add(time.Duration((i+1)*30) * time.Minute)).Format("Jan 2 15:04")))
			builder.WriteString(fmt.Sprintf("[EVENT] %s", timetable[i].event.Name))
		case timetable[i].deadline != nil:
			builder.WriteString(fmt.Sprintf("%s-%s: ", (currentTime.Add(time.Duration(i*30) * time.Minute)).Format("Jan 2 15:04"), (currentTime.Add(time.Duration((i+1)*30-5) * time.Minute)).Format("Jan 2 15:04")))
			builder.WriteString(fmt.Sprintf("[DEADLINE] %s", timetable[i].deadline.Name))
			builder.WriteString(fmt.Sprintln())
			builder.WriteString(fmt.Sprintf("%s-%s: 5 minute break", (currentTime.Add(time.Duration((i+1)*30-5) * time.Minute)).Format("Jan 2 15:04"), (currentTime.Add(time.Duration((i+1)*30) * time.Minute)).Format("Jan 2 15:04")))
		default:
			builder.WriteString(fmt.Sprintf("%s-%s: FREE SLOT", (currentTime.Add(time.Duration(i*30) * time.Minute)).Format("Jan 2 15:04"), (currentTime.Add(time.Duration((i+1)*30) * time.Minute)).Format("Jan 2 15:04")))
		}
		builder.WriteString(fmt.Sprintln())
	}
	return builder.String()
}
