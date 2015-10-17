package main

import (
	"github.com/jaffee/sked/scheduling"
	"strings"
)

type schedule struct {
	shifts []*shift
}

func NewSchedule(length int) *schedule {
	return &schedule{shifts: make([]*shift, length)}
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
