package types

import "time"

type Event struct {
	Name       string    `json:"name,omitempty" toml:"name"`
	StartTime  time.Time `json:"startTime" toml:"startTime"`
	EndTime    time.Time `json:"endTime" toml:"endTime"`
}

type Deadline struct {
	Name             string    `json:"name" toml:"name"`
	MinutesRemaining float64   `json:"minutesRemaining" toml:"minutesRemaining"`
	DeadlineTime         time.Time `json:"deadline" toml:"deadline"`
}