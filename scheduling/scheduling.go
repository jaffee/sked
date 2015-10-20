// Helper functions for building schedules
package scheduling

import (
	"time"
)

// A Shift represents a range of time. The shift starts at the time
// returned by Start (inclusive) and ends at the time returned by End
// (exclusive). This way a Shift that starts at the same time another
// one ends does not overlap, and there is no gap between the two
// Shifts.
type Shift interface {
	// Beginning of the shift - the shift includes this instant
	Start() time.Time

	// End of the shift - the shift does not include this instant
	End() time.Time

	// Return whether two shifts describe the same timespan. They do not
	// have to be equal byte-for-byte, just semantically
	// equivalent. (e.g. the times might be in different time zones)
	Equal(s2 Shift) bool

	// Return whether two shifts have any overlap
	Overlaps(s2 Shift) bool

	SetWorker(w Schedulable)

	// Return the Schedulable thing which is assigned to this shift (if any)
	Worker() Schedulable

	String() string
}

type Schedulable interface {
	// Returns whether this Schedulable is available for all of Shift s
	IsAvailable(s Shift) bool

	// State that this Schedulable is unavailable for the time period represented by Shift s
	AddUnavailable(s Shift)

	// An identifier for this schedulable (i.e. a person's name)
	Identifier() string

	// Priority for use in deciding when something should be scheduled
	Priority() int
	IncPriority(amnt int)
	DecPriority(amnt int)

	// Ordering to be used if priorities are the same
	Ordering() int
	SetOrdering(value int)
}

type Scheduleables interface {
	BuildSchedule() Schedule
}

type Schedule interface {
	Current() Schedulable
	String() string
}
