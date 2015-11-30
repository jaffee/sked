package main

import (
	"errors"
	"fmt"
	"time"
)

type Overlap int

const (
	Before        = Overlap(iota) // --{---}--[---]--  OR --{---}[---]--
	OverlapsStart                 // --{--[-}---]
	Prefix                        // --[{--}-]--
	EndsLater                     // --[{---]-}--
	Interior                      // --[-{-}-]--
	Same                          // --[{---}]--
	Subsumes                      // --{-[-]-}--
	StartsEarlier                 // --{-[--}]--
	Suffix                        // --[-{--}]--
	OverlapsEnd                   // --[--{-]--}--
	After                         // --[---]--{---}-- OR --[---]{---}--
)

// A Interval represents a range of time. The Interval starts at the time
// returned by Start (inclusive) and ends at the time returned by End
// (exclusive). This way a Interval that starts at the same time another
// one ends does not overlap, and there is no gap between the two
// Intervals.
type Intervaler interface {
	// Beginning of the Interval - the Interval includes this instant
	Start() time.Time

	SetStart(t time.Time) error

	// End of the Interval - the Interval does not include this instant
	End() time.Time

	SetEnd(t time.Time) error

	// Return whether two Intervals describe the same timespan. They do not
	// have to be equal byte-for-byte, just semantically
	// equivalent. (e.g. the times might be in different time zones)
	Equal(i2 Intervaler) bool

	// Return whether two Intervals have any overlap
	Overlaps(i2 Intervaler) bool

	// Return an Overlap which denotes how two Intervals overlap
	GetOverlap(i2 Intervaler) Overlap
}

// A Shift is an Interval with a worker (Schedulable) assigned to it.
type Shifter interface {
	// Beginning of the shift - the shift includes this instant
	Start() time.Time

	SetStart(t time.Time) error

	// End of the shift - the shift does not include this instant
	End() time.Time

	SetEnd(t time.Time) error

	// Return whether two shifts describe the same timespan. They do not
	// have to be equal byte-for-byte, just semantically
	// equivalent. (e.g. the times might be in different time zones)
	Equal(i2 Intervaler) bool

	// Return whether two shifts have any overlap
	Overlaps(i2 Intervaler) bool

	SetWorker(w *Person)

	// Return the Schedulable thing which is assigned to this shift (if any)
	Worker() *Person

	String() string

	// Return an Overlap which denotes how two Intervals overlap
	GetOverlap(i2 Intervaler) Overlap
}

type Interval struct {
	StartTime time.Time
	EndTime   time.Time
}

func NewInterval(start time.Time, end time.Time) (*Interval, error) {
	if !end.After(start) {
		return nil, errors.New("End must be after start.")
	}
	return &Interval{start, end}, nil
}

// Return the starting time of the shift
func (i *Interval) Start() time.Time {
	return i.StartTime
}

func (i *Interval) SetStart(t time.Time) error {
	if !t.Before(i.End()) {
		return errors.New("Start must be before end.")
	}
	i.StartTime = t
	return nil
}

// Return the end time of the shift
func (i *Interval) End() time.Time {
	return i.EndTime
}

func (i *Interval) SetEnd(t time.Time) error {
	if !t.After(i.Start()) {
		return errors.New("End must be after start.")
	}
	i.EndTime = t
	return nil
}

func (i *Interval) Equal(i2 Intervaler) bool {
	return i.Start().Equal(i2.Start()) && i.End().Equal(i2.End())
}

func (i *Interval) Overlaps(i2 Intervaler) bool {
	if !(i.End().Before(i2.Start()) || i.End().Equal(i2.Start())) {
		if !(i.Start().After(i2.End()) || i.Start().Equal(i2.End())) {
			return true
		}
	}
	return false
}

func (i *Interval) GetOverlap(i2 Intervaler) Overlap {
	if i.End().Before(i2.Start()) || i.End().Equal(i2.Start()) {
		return Before
	} else if i.Start().Before(i2.Start()) && i.End().Before(i2.End()) {
		return OverlapsStart
	} else if i.Start().Equal(i2.Start()) && i.End().Before(i2.End()) {
		return Prefix
	} else if i.Start().Equal(i2.Start()) && i.End().After(i2.End()) {
		return EndsLater
	} else if i.Start().After(i2.Start()) && i.End().Before(i2.End()) {
		return Interior
	} else if i.Equal(i2) {
		return Same
	} else if i.Start().Before(i2.Start()) && i.End().After(i2.End()) {
		return Subsumes
	} else if i.Start().Before(i2.Start()) && i.End().Equal(i2.End()) {
		return StartsEarlier
	} else if i.Start().After(i2.Start()) && i.End().Equal(i2.End()) {
		return Suffix
	} else if i.Start().Before(i2.End()) && i.End().After(i2.End()) {
		return OverlapsEnd
	} else if i.Start().After(i2.End()) || i.Start().Equal(i2.End()) {
		return After
	} else {
		panic(fmt.Sprintf("Bug in GetOverlap: i: %v, i2: %v\n", i, i2))
	}
}

type Shift struct {
	*Interval
	WorkerThing *Person
}

// Create a new Shift that goes from start to end.
func NewShift(start time.Time, end time.Time) (*Shift, error) {
	interval, err := NewInterval(start, end)
	if err != nil {
		return nil, err
	}
	ns := Shift{
		Interval:    interval,
		WorkerThing: NewPerson("EMPTY!"),
	}
	return &ns, nil
}

func (s *Shift) String() string {
	return fmt.Sprintf("%v from %v to %v", s.Worker().Identifier(), s.Start(), s.End())
}

func (s *Shift) Worker() *Person {
	return s.WorkerThing
}

func (s *Shift) SetWorker(w *Person) {
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
