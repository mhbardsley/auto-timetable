package main

import (
	"log"
	"time"
)

type timetableElement struct {
}

// function to generate the timetable
func generateTimetable(data inputData) {
	emptyTimetable := getEmptyTimetable(data.Deadlines)
	log.Println("Number of 30-minute segments in timetable: ", len(emptyTimetable))
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
