package main

import (
	"github.com/jaffee/sked/scheduling"
	"strings"
	"time"
)

type schedule struct {
	shifts []*shift
}

func NewSchedule(start time.Time, end time.Time, offset time.Weekday) *schedule {
	return &schedule{shifts: GetWeeklyShifts(start, end, offset)}
}

func (sched *schedule) Current() scheduling.Schedulable {
	return sched.shifts[0].Worker()
}

func (sched *schedule) String() string {
	sched_strings := make([]string, len(sched.shifts))
	for i, t := range sched.shifts {
		sched_strings[i] = t.String()
	}
	return strings.Join(sched_strings, "\n")

}

func (sched *schedule) Shifts() []*shift {
	return sched.shifts
}
