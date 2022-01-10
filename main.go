package main

import (
	"flag"
	"log"
	"math/rand"
	"time"
)

// global variable to store the current time
var currentTime time.Time = roundUp(time.Now())

func main() {
	rand.Seed(time.Now().Unix())
	filePtr := flag.String("f", "input.json", "The input's filename")
	inputData := getInput(filePtr)
	log.Println("Input data: ", inputData)

	generateTimetable(inputData)
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
