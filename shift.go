package main

import (
	"errors"
	"fmt"
	"github.com/jaffee/sked/scheduling"
	"time"
)

type Shift struct {
	StartTime   time.Time
	EndTime     time.Time
	WorkerThing scheduling.Schedulable
}

// Create a new Shift that goes from start to end.
func NewShift(start time.Time, end time.Time) (scheduling.Shift, error) {
	if end.Before(start) {
		return nil, errors.New("end must be after start")
	}
	ns := Shift{
		StartTime:   start,
		EndTime:     end,
		WorkerThing: NewPerson("EMPTY!"),
	}
	return &ns, nil
}

// Return the starting time of the shift
func (s *Shift) Start() time.Time {
	return s.StartTime
}

// Return the end time of the shift
func (s *Shift) End() time.Time {
	return s.EndTime
}

func (s *Shift) String() string {
	return fmt.Sprintf("%v from %v to %v", s.Worker().Identifier(), s.Start(), s.End())
}

func (s *Shift) Equal(s2 scheduling.Shift) bool {
	return s.Start().Equal(s2.Start()) && s.End().Equal(s2.End())
}

func (s *Shift) Overlaps(s2 scheduling.Shift) bool {
	if !(s.End().Before(s2.Start()) || s.End().Equal(s2.Start())) {
		if !(s.Start().After(s2.End()) || s.Start().Equal(s2.End())) {
			return true
		}
	}
	return false
}

func (s *Shift) Worker() scheduling.Schedulable {
	return s.WorkerThing
}

func (s *Shift) SetWorker(w scheduling.Schedulable) {
	s.WorkerThing = w
}

func GetWeeklyShifts(start time.Time, until time.Time, offset time.Weekday) []*Shift {
	lwd := getLastWeekday(start, offset)
	num_shifts := int((until.Sub(lwd).Hours()/24.0)/7.0) + 1
	shifts := make([]*Shift, num_shifts)
	cur := lwd
	for i := 0; i < num_shifts; i++ {
		startTime := cur
		cur = atMidnight(cur.Add(time.Hour * ((24 * 7) + 2)))
		ashift, err := NewShift(startTime, cur)
		if err != nil {
			panic(err)
		}

		shifts[i] = ashift
	}
	return shifts
}

func atMidnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// Get the beginning of the next day which is the day of the week
// denoted by `weekday` after or including `start`
func getLastWeekday(start time.Time, weekday time.Weekday) time.Time {
	cur := start
	for cur.Weekday() != weekday {
		cur = cur.Add(time.Hour * -23)
	}
	return atMidnight(cur)
}
