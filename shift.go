package main

import (
	"errors"
	"fmt"
	"github.com/jaffee/sked/scheduling"
	"time"
)

type shift struct {
	start  time.Time
	end    time.Time
	worker *person
}

// Create a new Shift that goes from start to end.
func NewShift(start time.Time, end time.Time) (scheduling.Shift, error) {
	if end.Before(start) {
		return nil, errors.New("end must be after start")
	}
	ns := shift{
		start: start,
		end:   end,
	}
	return &ns, nil
}

// Return the starting time of the shift
func (s *shift) Start() time.Time {
	return s.start
}

// Return the end time of the shift
func (s *shift) End() time.Time {
	return s.end
}

func (s *shift) Equal(s2 scheduling.Shift) bool {
	return s.Start().Equal(s2.Start()) && s.End().Equal(s2.End())
}

func (s *shift) Overlaps(s2 scheduling.Shift) bool {
	if !(s.End().Before(s2.Start()) || s.End().Equal(s2.Start())) {
		if !(s.Start().After(s2.End()) || s.Start().Equal(s2.End())) {
			return true
		}
	}
	return false
}

func (s *shift) Worker() scheduling.Schedulable {
	return s.worker
}

func (s *shift) SetWorker(w scheduling.Schedulable) {
	s.SetWorker(w)
}

func (s *shift) String() string {
	return fmt.Sprintf("%v: %v to %v", s.Worker().Identifier(), s.Start(), s.End())
}
