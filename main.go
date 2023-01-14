package main

import (
	"flag"
	"strconv"

	"github.com/mhbardsley/auto-timetable/backend"
)

func main() {
	filePtr := flag.String("f", "input.json", "The input's filename")
	slotsPtr := flag.String("s", "48", "The number of slots to display")
	repopulatePtr := flag.String("r", "0.04", "Repopulation threshold")

	flag.Parse()

	noOfSlots, _ := strconv.Atoi(*slotsPtr)
	threshold, _ := strconv.ParseFloat(*repopulatePtr, 64)
	inputData := backend.GetInput(filePtr, noOfSlots)

	backend.GenerateTimetable(inputData, threshold)
}
