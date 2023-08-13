package types

import "time"

type Event struct {
	Name       string    `json:"name,omitempty"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
}

type Deadline struct {
	Name             string    `json:"name"`
	MinutesRemaining float64   `json:"minutesRemaining"`
	Deadline         time.Time `json:"deadline"`
}