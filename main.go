package main

import (
	"flag"
	"math/rand"
	"strconv"
	"time"
)

// global variable to store the current time
var currentTime time.Time = roundUp(time.Now())

func main() {
	rand.Seed(currentTime.Unix())
	filePtr := flag.String("f", "input.json", "The input's filename")
	slotsPtr := flag.String("s", "48", "The number of slots to display")
	metaPtr := flag.String("m", "0.0625", "Chance (0-1) the app will ask to repopulate with deadlines and events")

	flag.Parse()

	noOfSlots, _ := strconv.Atoi(*slotsPtr)
	noOfMeta, _ := strconv.ParseFloat(*metaPtr, 64)
	inputData := getInput(filePtr, noOfSlots)

	generateTimetable(inputData, noOfMeta)
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
