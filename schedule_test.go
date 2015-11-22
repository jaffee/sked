package main

import (
	"testing"
	"time"
)

func TestNext(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	start := time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	end := time.Date(2015, time.October, 18, 0, 0, 0, 0, loc)

	sched := NewSchedule(start, end, time.Wednesday)

	numShifts := sched.NumShifts()

	for i := 0; i <= numShifts; i++ {
		_, err := sched.Next()
		if i == numShifts {
			if err == nil {
				t.Fatalf("Error should not be nil. Sched is %v", sched)
			}
		} else { // i < numShifts
			if err != nil {
				t.Fatalf("Shouldn't be an err, Sched %v", sched)
			}
		}
	}

}
