package backend

import (
	"math/rand"
	"time"
)

// global variable to store the current time
var currentTime time.Time = roundUp(time.Now())

func init() {
	rand.Seed(currentTime.Unix())
}