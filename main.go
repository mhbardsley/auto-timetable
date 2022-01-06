package main

import (
	"flag"
	"log"
	"time"
)

// global variable to store the current time
var currentTime time.Time = time.Now()

func main() {
	filePtr := flag.String("f", "input.json", "The input's filename")
	inputData := getInput(filePtr)

	log.Println("Input data: ", inputData)

	generateTimetable(inputData)
}
