package main

import (
	"errors"
	"fmt"
	"github.com/jaffee/sked/scheduling"
	"sort"
	"time"
)

type state struct {
	people   map[string]*person
	offset   time.Weekday
	schedule *schedule
}

func (s *state) BuildSchedule(start time.Time, end time.Time) scheduling.Schedule {
	sched := NewSchedule(start, end, s.offset)
	personList := tempPersonList(s.people)
	for _, cur_shift := range sched.Shifts() {

		// find person with lowest priority who is available
		np, err := nextAvailable(personList, cur_shift)
		if err != nil {
			np = &person{name: "EMPTY!"}
		}

		cur_shift.SetWorker(s.people[np.name])

		// re-calc priorities
		fmt.Println("before loop")
		fmt.Println(personList)
		for _, p := range personList {
			if p.Identifier() != np.Identifier() {
				p.DecPriority(1)
			} else {
				p.IncPriority(len(personList))
			}
			fmt.Println(p, p.Priority())
		}
		fmt.Println("after loop")
		fmt.Println(personList)
	}
	return sched
}

type ByPriority []*person

func (b ByPriority) Len() int      { return len(b) }
func (b ByPriority) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByPriority) Less(i, j int) bool {
	if b[i].Priority() == b[j].Priority() {
		if b[i].Ordering() == b[j].Ordering() {
			return b[i].Identifier() < b[j].Identifier()
		} else {
			return b[i].Ordering() < b[j].Ordering()
		}
	} else {
		return b[i].Priority() < b[j].Priority()
	}
}

func nextAvailable(personList []*person, cur_shift scheduling.Shift) (*person, error) {
	sort.Sort(ByPriority(personList))
	var np *person
	found := false
	for _, p := range personList {
		if p.IsAvailable(cur_shift) {
			np = p
			found = true
			break
		}
	}
	if found {
		return np, nil
	} else {
		return &person{}, errors.New("Could not find anyone to work the shift")
	}
}

func tempPersonList(people map[string]*person) []*person {
	personList := make([]*person, len(people))
	i := 0
	for _, p := range people {
		personList[i] = &person{}
		personList[i].name = p.name
		personList[i].unavailability = p.unavailability
		personList[i].priority = p.priority
		personList[i].orderNum = p.orderNum
		i += 1
	}
	return personList
}

func NewState(offset time.Weekday) state {
	// Wednesday is the default for offset... makes sense right?
	s := state{
		people: make(map[string]*person),
		offset: offset,
	}
	return s
}

// func checkSchedule(s state) {
// 	if len(s.schedule) == 0 {
// 		return
// 	}
// 	now := time.Now()
// 	if now.After(s.schedule[0].End()) {
// 		curWorker = s.schedule[0].Worker()
// 		for n, p := range s.people {
// 			if p.Identifier() != curWorker.Identifier() {
// 				p.DecPriority(1)
// 			} else {
// 				p.IncPriority(len(s.people))
// 			}
// 		}
// 	}
// }
