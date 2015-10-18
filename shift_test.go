package main

import (
	"testing"
	"time"
)

func TestNewShiftStartAndEnd(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	start := time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	end := time.Date(2015, time.October, 12, 0, 0, 0, 0, loc)

	s, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Problem creating new shift - err should be nil")
	}

	if !s.Start().Equal(start) {
		t.Fatalf("Shift starting time got mucked up somehow")
	}

	if !s.End().Equal(end) {
		t.Fatalf("Shift ending time got mucked up somehow")
	}

	_, err = NewShift(end, start)
	if err == nil {
		t.Fatalf("Problem creating new shift - err should not be nil because start was before end")
	}

}

func TestEqual(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	start := time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	end := time.Date(2015, time.October, 12, 0, 0, 0, 0, loc)
	chicShift, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Problem creating new shift: %v", err)
	}

	loc, err = time.LoadLocation("America/New_York")

	start = time.Date(2015, time.October, 10, 1, 0, 0, 0, loc)
	end = time.Date(2015, time.October, 12, 1, 0, 0, 0, loc)
	nycShift, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Problem creating new shift: %v", err)
	}

	if !chicShift.Equal(nycShift) {
		t.Fatalf("Shifts should be equal even with different timezones")
	}
}
