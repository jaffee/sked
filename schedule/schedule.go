// Helper functions for building schedules
package schedule

import (
	"fmt"
	"time"
)

type Shift struct {
	Start myTime
	End   myTime
}

func (s Shift) Equal(s2 Shift) bool {
	return s.Start.Equal(s2.Start.Time) && s.End.Equal(s2.End.Time)
}

func (s Shift) String() string {
	return fmt.Sprintf("%v to %v\n", s.Start.String(), s.End.String())
}

type myTime struct {
	time.Time
}

func (t myTime) atMidnight() myTime {
	return myTime{time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())}
}

func (t myTime) String() string {
	return t.Time.Format(time.UnixDate)
}

func getWeeklyShifts(start time.Time, until time.Time, offset time.Weekday) []Shift {
	lwd := getLastWeekday(start, offset)
	num_shifts := int((until.Sub(lwd.Time).Hours()/24.0)/7.0) + 1
	shifts := make([]Shift, num_shifts)
	cur := lwd
	for i := 0; i < num_shifts; i++ {
		shifts[i].Start = cur
		cur = myTime{cur.Add(time.Hour * ((24 * 7) + 2))}.atMidnight()
		shifts[i].End = cur
	}
	return shifts
}

// Get the beginning of the next day which is the day of the week
// denoted by `weekday` after or including `start`
func getLastWeekday(start time.Time, weekday time.Weekday) myTime {
	cur := myTime{start}
	for cur.Weekday() != weekday {
		cur = myTime{cur.Add(time.Hour * -23)}
	}
	return cur.atMidnight()
}
