package main

import (
	"math"
	"math/rand"
)

// fill the timetable with deadlines probabilistically
// assume that it is possible
func fillTimetable(timetable []timetableElement, deadlines []deadline) {
	for i := 0; ; i++ {
		// need to construct a slice of weights
		if hasFilled(timetable, deadlines, i) {
			break
		}
	}
}

// hasFilled will check if deadlines have been satisfied with a power of pow
func hasFilled(timetable []timetableElement, deadlines []deadline, pow int) bool {
	var chosenIndex int
	deadlinesCopy := copyDeadlines(deadlines)
	for i, slot := range timetable {
		if slot.event == nil && len(deadlinesCopy) > 0 {
			weights := getWeights(deadlinesCopy, pow)
			cumulateWeights(weights)
			r := rand.Float64() * weights[len(weights)-1]
			for j, weight := range weights {
				if r <= weight {
					chosenIndex = j
					break
				}
			}
			timetable[i].deadline = &(deadlinesCopy[chosenIndex])
			deadlinesCopy = reduceDeadlines(deadlinesCopy, chosenIndex)
			if _, possible := possibleTimetabling(deadlinesCopy); !possible {
				return false
			}
		}
	}
	copy(deadlines, deadlinesCopy)
	return true
}

func copyDeadlines(deadlines []deadline) (deadlinesCopy []deadline) {
	deadlinesCopy = make([]deadline, len(deadlines))
	_ = copy(deadlinesCopy, deadlines)
	return deadlinesCopy
}

// getWeights goes over deadlines and makes a deciated float slice
func getWeights(deadlines []deadline, pow int) (weights []float64) {
	weights = make([]float64, len(deadlines))
	for i, deadline := range deadlines {
		weights[i] = math.Pow(float64(deadline.slotsRemaining)/float64(deadline.slotsAvailable), float64(pow))
	}
	return weights
}

// cumulateWeights adds the weights as they go along
func cumulateWeights(weights []float64) {
	for i, weight := range weights {
		if i == 0 {
			weights[i] = weight
		} else {
			weights[i] = weight + weights[i-1]
		}
	}
}

// change the deadlines so that they will delete if complete, reduce otherwise
func reduceDeadlines(deadlines []deadline, index int) []deadline {
	zeroFlag := false
	for i, deadline := range deadlines {
		deadlines[i].slotsAvailable = deadline.slotsAvailable - 1
		if i == index {
			deadlines[i].slotsRemaining = deadline.slotsRemaining - 1
			if deadlines[i].slotsRemaining <= 0 {
				zeroFlag = true
			}
		}
	}
	if zeroFlag {
		newDeadlines := make([]deadline, len(deadlines)-1)
		copy(newDeadlines[:index], deadlines[:index])
		copy(newDeadlines[index:], deadlines[index+1:])
		return newDeadlines
	}
	return deadlines
}
