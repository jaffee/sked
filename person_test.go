package main

import (
	"testing"
	"time"
)

func TestAddUnavailable(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	start := time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	end := time.Date(2015, time.October, 12, 0, 0, 0, 0, loc)

	p := NewPerson("joe")
	aShift, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Failed to make a shift: %v", err)
	}
	p.AddUnavailable(aShift)

	if len(p.unavailability) != 1 {
		t.Fatalf("p should have one unavailability, not %v", len(p.unavailability))
	}

	start = time.Date(2015, time.October, 22, 0, 0, 0, 0, loc)
	end = time.Date(2015, time.October, 23, 0, 0, 0, 0, loc)
	aShift, err = NewShift(start, end)
	if err != nil {
		t.Fatalf("Failed to make a shift: %v", err)
	}

	p.AddUnavailable(aShift)

	if len(p.unavailability) != 2 {
		t.Fatalf("p should have two unavailabilities, not %v", len(p.unavailability))
	}

}

func TestIsAvailable(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	start := time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	end := time.Date(2015, time.October, 12, 0, 0, 0, 0, loc)

	p := NewPerson("joe")
	aShift, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Failed to make a shift: %v", err)
	}
	p.AddUnavailable(aShift)

	// make previous shift
	start = time.Date(2015, time.October, 8, 0, 0, 0, 0, loc)
	end = time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	prevShift, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Failed to make a shift: %v", err)
	}
	if !p.IsAvailable(prevShift) {
		t.Fatalf("Should be available for %v", prevShift)
	}

	// make after shift
	start = time.Date(2015, time.October, 12, 0, 0, 0, 0, loc)
	end = time.Date(2015, time.October, 15, 0, 0, 0, 0, loc)
	afterShift, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Failed to make a shift: %v", err)
	}
	if !p.IsAvailable(afterShift) {
		t.Fatalf("Should be available for %v", afterShift)
	}

	// make encompassing shift
	start = time.Date(2015, time.October, 8, 0, 0, 0, 0, loc)
	end = time.Date(2015, time.October, 15, 0, 0, 0, 0, loc)
	encompassingShift, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Failed to make a shift: %v", err)
	}
	if p.IsAvailable(encompassingShift) {
		t.Fatalf("Should not be available for %v", encompassingShift)
	}

	// make overlapping shift
	start = time.Date(2015, time.October, 9, 0, 0, 0, 0, loc)
	end = time.Date(2015, time.October, 11, 0, 0, 0, 0, loc)
	overlappingShift, err := NewShift(start, end)
	if err != nil {
		t.Fatalf("Failed to make a shift: %v", err)
	}
	if p.IsAvailable(overlappingShift) {
		t.Fatalf("Should not be available for %v", overlappingShift)
	}
}

func TestIdentifier(t *testing.T) {
	p := NewPerson("joe")
	if p.Identifier() != "joe" {
		t.Fatalf("Identifier should be joe")
	}
}

func TestPriority(t *testing.T) {
	startingPriority := 0
	p := NewPerson("joe")
	if p.Priority() != startingPriority {
		t.Fatalf("Initial priority should be %v", startingPriority)
	}

	p.IncPriority(3)
	expectedPriority := startingPriority + 3
	actual := p.Priority()
	if actual != expectedPriority {
		t.Fatalf("Problem increasing priority, expected: %v, got: %v", expectedPriority, actual)
	}

	p.DecPriority(9)
	expectedPriority = actual - 9
	actual = p.Priority()
	if actual != expectedPriority {
		t.Fatalf("Problem decreasing priority, expected: %v, got: %v", expectedPriority, actual)
	}

}
