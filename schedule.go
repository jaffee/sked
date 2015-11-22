package main

import (
	"errors"
	"fmt"
	"github.com/jaffee/sked/scheduling"
	"strings"
	"time"
)

const MAX_WEEKS = 10

type Schedule struct {
	ShiftsList []*Shift
	shiftIdx   int
}

func NewSchedule(start time.Time, end time.Time, offset time.Weekday) *Schedule {
	return &Schedule{ShiftsList: GetWeeklyShifts(start, end, offset)}
}

func (sched *Schedule) Current() scheduling.Schedulable {
	return sched.ShiftsList[0].Worker()
}

func (sched *Schedule) String() string {
	sched_strings := make([]string, len(sched.ShiftsList))
	for i, t := range sched.ShiftsList {
		sched_strings[i] = t.String()
	}
	return strings.Join(sched_strings, "\n")

}

func (sched *Schedule) Next() (scheduling.Shift, error) {
	if sched.shiftIdx < len(sched.ShiftsList) {
		sched.shiftIdx += 1
		return sched.ShiftsList[sched.shiftIdx-1], nil
	} else {
		sched.shiftIdx = 0
		return &Shift{}, errors.New("End of iteration")
	}
}

func (sched *Schedule) NumShifts() int {
	return len(sched.ShiftsList)
}

func (s *Schedule) SPrintCalendar() string {
	startShifts := s.ShiftsList[0].Start()
	start := getLastWeekday(startShifts, time.Sunday)

	line := "| Sunday    | Monday    | Tuesday   | Wednesday | Thursday  | Friday    | Saturday  |\n"
	line += "|-----------+-----------+-----------+-----------+-----------+-----------+-----------|\n"
	shiftIdx := 0
	for i := 0; i < MAX_WEEKS; i++ {
		line += "|"
		for j := 0; j < 7; j++ {
			curShift := s.ShiftsList[shiftIdx]
			dayNum := start.Day() + (7 * i) + j
			day := time.Date(start.Year(), start.Month(), dayNum,
				start.Hour(), start.Minute(), start.Second(),
				start.Nanosecond(), start.Location())
			if day.Before(curShift.Start()) {
				line += "           |"
			} else {
				if day.After(curShift.End()) {
					shiftIdx += 1
					if shiftIdx >= len(s.ShiftsList) {
						break
					} else {
						curShift = s.ShiftsList[shiftIdx]
					}
				}
				line += fmt.Sprintf(" %2v %-6v |", day.Day(), curShift.Worker().Identifier())
			}
		}
		line += "\n"
	}
	return line
}
