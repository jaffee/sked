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
