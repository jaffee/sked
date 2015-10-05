// Helper functions for building schedules
package schedule

import (
	"fmt"
	"time"
)

type Shift interface {
	Start() time.Time
	End() time.Time
}

type shift struct {
	start time.Time
	end   time.Time
}

func (s shift) Start() time.Time {
	return s.start
}

func (s shift) End() time.Time {
	return s.end
}

func (s shift) Equal(s2 shift) bool {
	return s.start.Equal(s2.start) && s.end.Equal(s2.end)
}

func (s shift) String() string {
	return fmt.Sprintf("%v to %v\n", s.start.String(), s.end.String())
}

func atMidnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetWeeklyShifts(start time.Time, until time.Time, offset time.Weekday) []Shift {
	lwd := getLastWeekday(start, offset)
	num_shifts := int((until.Sub(lwd).Hours()/24.0)/7.0) + 1
	shifts := make([]Shift, num_shifts)
	cur := lwd
	var ashift shift
	for i := 0; i < num_shifts; i++ {
		ashift = shift{}
		ashift.start = cur
		cur = atMidnight(cur.Add(time.Hour * ((24 * 7) + 2)))
		ashift.end = cur
		shifts[i] = ashift
	}
	return []Shift(shifts)
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
