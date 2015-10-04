// Helper functions for building schedules
package schedule

import (
	"fmt"
	"time"
)

type Shift struct {
	Start time.Time
	End   time.Time
}

func (s Shift) Equal(s2 Shift) bool {
	return s.Start.Equal(s2.Start) && s.End.Equal(s2.End)
}

func (s Shift) String() string {
	return fmt.Sprintf("%v to %v\n", s.Start.String(), s.End.String())
}

func atMidnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func getWeeklyShifts(start time.Time, until time.Time, offset time.Weekday) []Shift {
	lwd := getLastWeekday(start, offset)
	num_shifts := int((until.Sub(lwd).Hours()/24.0)/7.0) + 1
	shifts := make([]Shift, num_shifts)
	cur := lwd
	for i := 0; i < num_shifts; i++ {
		shifts[i].Start = cur
		cur = atMidnight(cur.Add(time.Hour * ((24 * 7) + 2)))
		shifts[i].End = cur
	}
	return shifts
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
