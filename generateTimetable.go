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
	roundedT1 := roundUp(time1)
	roundedT2 := roundDown(time2)

	durationBetween := roundedT2.Sub(roundedT1)
	return (int)(durationBetween.Minutes() / 30)
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
