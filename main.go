package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/mhbardsley/auto-timetable/backend"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: auto-timetable <subcommand> [flags]")
		fmt.Println("subcommands:")
		fmt.Println("  generate - generate a timetable")
		fmt.Println("  help - display this help")
	}
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
		filePtr := generateCmd.String("f", "input.json", "The input's filename")
		slotsPtr := generateCmd.String("s", "48", "The number of slots to display")
		repopulatePtr := generateCmd.String("r", "0.04", "Repopulation threshold")
		generateCmd.Parse(args[1:])
		noOfSlots, _ := strconv.Atoi(*slotsPtr)
		threshold, _ := strconv.ParseFloat(*repopulatePtr, 64)
		inputData := backend.GetInput(filePtr, noOfSlots)
		backend.GenerateTimetable(inputData, threshold)
	case "help":
		flag.Usage()
	default:
		fmt.Printf("%q is not valid subcommand.\n", args[0])
		flag.Usage()
		os.Exit(1)
	}
}
