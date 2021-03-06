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

	start = time.Date(2015, time.October, 10, 2, 0, 0, 0, loc)
	end = time.Date(2015, time.October, 12, 2, 0, 0, 0, loc)
	nycShift, err = NewShift(start, end)
	if err != nil {
		t.Fatalf("Problem creating new shift: %v", err)
	}

	if chicShift.Equal(nycShift) {
		t.Fatalf("Shifts should not equal because they represent different times")
	}
}

func TestOverlaps(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	start := time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	end := time.Date(2015, time.October, 12, 0, 0, 0, 0, loc)

	s, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Problem creating new shift - err should be nil")
	}

	// same times
	start = time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	end = time.Date(2015, time.October, 12, 0, 0, 0, 0, loc)

	sameShift, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Problem creating new shift - err should be nil")
	}

	if !s.Overlaps(sameShift) {
		t.Fatalf("Two shifts representing the same time should overlap")
	}

	// butting up
	start = time.Date(2015, time.October, 8, 0, 0, 0, 0, loc)
	end = time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)

	beforeShift, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Problem creating new shift - err should be nil")
	}

	if s.Overlaps(beforeShift) {
		t.Fatalf("Two shifts right next to each other should not overlap")
	}

	// partial overlap
	start = time.Date(2015, time.October, 8, 0, 0, 0, 0, loc)
	end = time.Date(2015, time.October, 11, 0, 0, 0, 0, loc)

	partialShift, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Problem creating new shift - err should be nil")
	}

	if !s.Overlaps(partialShift) {
		t.Fatalf("Shifts that partially overlap should overlap")
	}
}

func TestWorkers(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	start := time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	end := time.Date(2015, time.October, 12, 0, 0, 0, 0, loc)
	s, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Problem creating new shift - err should be nil")
	}

	p := NewPerson("joe")
	s.SetWorker(p)

	if s.Worker() != p {
		t.Fatalf("Worker should be the worker...")
	}

}

func TestGetWeeklyShifts(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	start := time.Date(2015, time.October, 11, 22, 3, 0, 0, loc)
	until := time.Date(2015, time.November, 5, 14, 1, 0, 0, loc)
	shifts := GetWeeklyShifts(start, until, time.Wednesday)
	if len(shifts) != 5 {
		t.Fatalf("Wrong number of shifts: %v", shifts)
	}

	expectedStrings := []string{
		"Wed Oct  7 00:00:00 CDT 2015",
		"Wed Oct 14 00:00:00 CDT 2015",
		"Wed Oct 14 00:00:00 CDT 2015",
		"Wed Oct 21 00:00:00 CDT 2015",
		"Wed Oct 21 00:00:00 CDT 2015",
		"Wed Oct 28 00:00:00 CDT 2015",
		"Wed Oct 28 00:00:00 CDT 2015",
		"Wed Nov  4 00:00:00 CST 2015",
		"Wed Nov  4 00:00:00 CST 2015",
		"Wed Nov 11 00:00:00 CST 2015",
	}
	expectedShifts := make([]*Shift, 5)
	for i := 0; i < 5; i++ {
		dStart, err := time.Parse(time.UnixDate, expectedStrings[2*i])
		if err != nil {
			t.Fatalf("error parsing date for expected")
		}
		dEnd, err := time.Parse(time.UnixDate, expectedStrings[2*i+1])
		if err != nil {
			t.Fatalf("error parsing date for expected")
		}
		s, err := NewShift(dStart, dEnd)
		if err != nil {
			t.Fatalf("Error creating shift: %v", err)
		}
		expectedShifts[i] = s
	}
	for i, s := range shifts {
		if !expectedShifts[i].Equal(s) {
			t.Fatalf("Output failed to match at index:%v, \nexpected:%vactual:  %v", i, expectedShifts[i], s)
		}
	}
}

func TestGetOverlap(t *testing.T) {
	var istart, iend time.Time
	var i Interval
	loc, _ := time.LoadLocation("America/Chicago")
	i2start := time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	i2end := time.Date(2015, time.October, 13, 0, 0, 0, 0, loc)
	i2 := &Interval{i2start, i2end}

	istart = time.Date(2015, time.October, 5, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 8, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != Before {
		t.Fatalf("GetOverlap Before 1")
	}

	istart = time.Date(2015, time.October, 5, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != Before {
		t.Fatalf("GetOverlap Before 2")
	}

	istart = time.Date(2015, time.October, 5, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 11, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != OverlapsStart {
		t.Fatalf("GetOverlap OverlapsStart 3")
	}

	istart = time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 11, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != Prefix {
		t.Fatalf("GetOverlap Prefix 4")
	}

	istart = time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 14, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != EndsLater {
		t.Fatalf("GetOverlap EndsLater 5")
	}

	istart = time.Date(2015, time.October, 11, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 12, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != Interior {
		t.Fatalf("GetOverlap Interior 6")
	}

	istart = time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 13, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != Same {
		t.Fatalf("GetOverlap Same 7")
	}

	istart = time.Date(2015, time.October, 9, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 14, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != Subsumes {
		t.Fatalf("GetOverlap Subsumes 8")
	}

	istart = time.Date(2015, time.October, 9, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 13, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != StartsEarlier {
		t.Fatalf("GetOverlap StartsEarlier 9")
	}

	istart = time.Date(2015, time.October, 11, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 13, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != Suffix {
		t.Fatalf("GetOverlap Suffix 10")
	}

	istart = time.Date(2015, time.October, 11, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 14, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != OverlapsEnd {
		t.Fatalf("GetOverlap OverlapsEnd 11")
	}

	istart = time.Date(2015, time.October, 14, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 15, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != After {
		t.Fatalf("GetOverlap After 12")
	}

	istart = time.Date(2015, time.October, 13, 0, 0, 0, 0, loc)
	iend = time.Date(2015, time.October, 15, 0, 0, 0, 0, loc)
	i = Interval{istart, iend}
	if i.GetOverlap(i2) != After {
		t.Fatalf("GetOverlap After 13")
	}

}
