package main

import (
	"time"
)

func runSchedule(skedState *State) {
	// lastRun := time.Date(0, 1, 1, 0, 0, 0, 0, time.Local)
	for {
		skedState.Lock()

		skedState.Unlock()
		// lastRun := time.Now()
		time.Sleep(time.Second * 10)
	}
}
