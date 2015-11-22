package main

import (
	"testing"
	"time"
)

func TestPersist(t *testing.T) {
	s := NewState(time.Wednesday)
	loc, _ := time.LoadLocation("America/Chicago")
	start := time.Date(2015, time.October, 11, 22, 0, 0, 0, loc)
	end := time.Date(2015, time.November, 18, 17, 0, 0, 0, loc)
	s.BuildSchedule(start, end)
	s.Persist()
}
