package main

import (
	"flag"
	"log"
)

func main() {
	filePtr := flag.String("f", "input.json", "The input's filename")
	inputData := getInput(filePtr)

	log.Println(inputData)
}
