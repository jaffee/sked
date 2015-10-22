package main

import (
	"github.com/jaffee/sked/scheduling"
	"strings"
	"time"
)

type Schedule struct {
	ShiftsList []*Shift
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

func (sched *Schedule) Shifts() []*Shift {
	return sched.ShiftsList
}
