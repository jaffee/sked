package scheduling

import (
	"testing"
	"time"
)

func TestGetLastWeekday(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	start := time.Date(2015, time.October, 11, 22, 3, 0, 0, loc)
	lwd := getLastWeekday(start, time.Wednesday)
	if lwd != time.Date(2015, time.October, 7, 0, 0, 0, 0, loc) {
		t.Fatalf("Last Wednesday came back as '%v'", lwd)
	}
}

func TestGetLastWeekdaySameDay(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	start := time.Date(2015, time.October, 7, 22, 3, 0, 0, loc)
	lwd := getLastWeekday(start, time.Wednesday)
	if lwd != time.Date(2015, time.October, 7, 0, 0, 0, 0, loc) {
		t.Fatalf("Last Wednesday came back as '%v'", lwd)
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
	expectedShifts := make([]shift, 5)
	for i := 0; i < 5; i++ {
		dStart, err := time.Parse(time.UnixDate, expectedStrings[2*i])
		if err != nil {
			t.Fatalf("error parsing date for expected")
		}
		dEnd, err := time.Parse(time.UnixDate, expectedStrings[2*i+1])
		if err != nil {
			t.Fatalf("error parsing date for expected")
		}
		expectedShifts[i].start = dStart
		expectedShifts[i].end = dEnd
	}
	for i, s := range shifts {
		if !expectedShifts[i].Equal(s) {
			t.Fatalf("Output failed to match at index:%v, \nexpected:%vactual:  %v", i, expectedShifts[i], s)
		}
	}
}

func TestNewShift(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	s := NewShift(time.Date(2015, time.April, 5, 0, 0, 0, 0, loc), time.Date(2015, time.April, 8, 0, 0, 0, 0, loc))
	if sta := s.Start(); sta != time.Date(2015, time.April, 5, 0, 0, 0, 0, loc) {
		t.Fatalf("Unexpected behavior from Shift - s.Start: %v", sta)
	}
	if end := s.End(); end != time.Date(2015, time.April, 8, 0, 0, 0, 0, loc) {
		t.Fatalf("Unexpected behavior from Shift - s.End: %v", end)
	}
}

func TestOverlaps(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	s1 := NewShift(time.Date(2015, time.April, 5, 0, 0, 0, 0, loc), time.Date(2015, time.April, 8, 0, 0, 0, 0, loc))
	s2 := NewShift(time.Date(2015, time.April, 2, 0, 0, 0, 0, loc), time.Date(2015, time.April, 9, 0, 0, 0, 0, loc))
	if !s1.Overlaps(s2) {
		t.Fatalf("s2 completely covers s1 - these should overlap")
	}
	s1 = NewShift(time.Date(2015, time.April, 5, 0, 0, 0, 0, loc), time.Date(2015, time.April, 8, 0, 0, 0, 0, loc))
	s2 = NewShift(time.Date(2015, time.April, 2, 0, 0, 0, 0, loc), time.Date(2015, time.April, 5, 0, 0, 0, 0, loc))
	if s1.Overlaps(s2) {
		t.Fatalf("s2 bumps up against s1 - these should not overlap")
	}
}