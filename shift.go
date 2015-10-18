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
	worker scheduling.Schedulable
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
	s.worker = w
}

func (s *shift) String() string {
	return fmt.Sprintf("%v: %v to %v", s.Worker().Identifier(), s.Start(), s.End())
}

func GetWeeklyShifts(start time.Time, until time.Time, offset time.Weekday) []*shift {
	lwd := getLastWeekday(start, offset)
	num_shifts := int((until.Sub(lwd).Hours()/24.0)/7.0) + 1
	shifts := make([]*shift, num_shifts)
	cur := lwd
	var ashift *shift
	for i := 0; i < num_shifts; i++ {
		ashift = &shift{}
		ashift.start = cur
		cur = atMidnight(cur.Add(time.Hour * ((24 * 7) + 2)))
		ashift.end = cur
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
