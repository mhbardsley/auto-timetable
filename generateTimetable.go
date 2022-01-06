package main

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

type timetableElement struct {
	event    *string
	deadline *string
}

// function to generate the timetable
func generateTimetable(data inputData) {
	timetable := getEmptyTimetable(data.Deadlines)
	log.Println("Number of 30-minute segments in timetable: ", len(timetable))

	fillWithEvents(timetable, data.Events)
	log.Println("Timetable filled with events:")
	log.Printf("%s", printTimetable(timetable))
}

// generate a slice of timetable elements
func getEmptyTimetable(deadlines []deadline) (timetable []timetableElement) {
	var numberOfSpaces int
	noOfDeadlines := len(deadlines)
	if noOfDeadlines == 0 {
		numberOfSpaces = 0
	} else {
		numberOfSpaces = segmentsBetween(currentTime, deadlines[noOfDeadlines-1].Deadline)
	}
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
		startIndex = segmentsBetween(currentTime, event.StartTime)
		endIndex = int(math.Min(float64(segmentsBetween(currentTime, event.EndTime)), float64(len(timetable))))
		selectedElements = timetable[startIndex:endIndex]
		for j := range selectedElements {
			selectedElements[j].event = &(events[i].Name)
			log.Println("Event added to timetable: ", event.Name, " with address", selectedElements[j].event)
		}
	}
}

// function to print the timetable as-is
func printTimetable(timetable []timetableElement) string {
	builder := strings.Builder{}
	for i, slot := range timetable {
		builder.WriteString(fmt.Sprintf("%s-%s: ", (currentTime.Add(time.Duration(i*30) * time.Minute)).Format("Jan 2 15:04"), (currentTime.Add(time.Duration((i+1)*30) * time.Minute)).Format("Jan 2 15:04")))
		if slot.event != nil {
			log.Println("Event pointed to has address: ", slot.event)
			builder.WriteString(fmt.Sprintf("[EVENT] %s", *(slot.event)))
		} else if slot.deadline != nil {
			builder.WriteString(fmt.Sprintf("[DEADLINE] %s", *(slot.deadline)))
		}
		builder.WriteString(fmt.Sprintln())
	}
	return builder.String()
}
