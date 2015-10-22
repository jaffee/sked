package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/jaffee/sked/scheduling"
	"os"
	"sort"
	"time"
)

type State struct {
	People    map[string]*Person
	Offset    time.Weekday
	Schedule  *Schedule
	StorageID string
}

func (s *State) Persist() error {
	f, err := os.Create(s.StorageID)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	err = enc.Encode(s)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (s *State) Populate() error {
	f, err := os.Open(s.StorageID)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	dec := gob.NewDecoder(r)
	err = dec.Decode(s)
	if err != nil {
		return err
	}
	return nil
}

func (s *State) BuildSchedule(start time.Time, end time.Time) scheduling.Schedule {
	sched := NewSchedule(start, end, s.Offset)
	personList := tempPersonList(s.People)
	for _, cur_shift := range sched.Shifts() {

		// find person with lowest priority who is available
		np, err := nextAvailable(personList, cur_shift)
		if err != nil {
			np = &Person{Name: "EMPTY!"}
		}

		cur_shift.SetWorker(s.People[np.Name])

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

func nextAvailable(personList []*Person, cur_shift scheduling.Shift) (*Person, error) {
	sort.Sort(ByPriority(personList))
	var np *Person
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
		return &Person{}, errors.New("Could not find anyone to work the shift")
	}
}

func tempPersonList(people map[string]*Person) []*Person {
	personList := make([]*Person, len(people))
	i := 0
	for _, p := range people {
		personList[i] = &Person{}
		personList[i].Name = p.Name
		personList[i].Unavailability = p.Unavailability
		personList[i].PriorityNum = p.PriorityNum
		personList[i].OrderNum = p.OrderNum
		i += 1
	}
	return personList
}

func NewState(offset time.Weekday) State {
	// Wednesday is the default for offset... makes sense right?
	s := State{
		People:    make(map[string]*Person),
		Offset:    offset,
		StorageID: "skedState.gob",
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
