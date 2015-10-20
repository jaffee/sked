/*

sked - Slackbot for round robin scheduling

Based on mybot - an illustrative slackbot in Go, copyright reproduced below.

Copyright (c) 2015 RapidLoop

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/jaffee/sked/scheduling"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type state struct {
	people map[string]*person
	offset time.Weekday
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
		return b[i].Identifier() < b[j].Identifier()
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

type command struct {
	action string
	args   []string
}

type action struct {
	function func(command, *state) string
	help     string
}

func writeHandler(comChan chan string, w *bufio.Writer) {
	for {
		s := <-comChan + "\n"
		fmt.Println("Writing", s, []byte(s))
		n, err := w.Write([]byte(s))
		w.Flush()
		fmt.Println(n, "bytes written")
		if err != nil {
			log.Printf("Problem while writing, err:%v", err)
			// TODO - set up some kind of recovery, or at least notify slack of failure
			panic(err)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: sked <slack-bot-token> [state-file]\n")
		os.Exit(1)
	}

	// set up state
	command_map := map[string]action{
		"current":  action{getCurrent, "Tell me who's scheduled right now"},
		"add":      action{addPerson, "Add a new person to be scheduled"},
		"list":     action{list, "List all the possible people that could be scheduled"},
		"unavail":  action{addUnavailable, "unavail <name> <[YYYY]MMDD[HH]> [to [YYYY]MMDD[HH]]"},
		"schedule": action{getSchedule, "Build the schedule using the people and availabilities given so far"},
	}
	sked_state := NewState(time.Wednesday)

	// Output file handling
	var filename string
	var stateFile *os.File
	if len(os.Args) >= 3 {
		filename = os.Args[2]
	} else {
		filename = "sked-log.txt"
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		stateFile, err = os.Create(filename)
		if err != nil {
			log.Fatalf("Could not open file: %v, error: %v", filename, err)
		}
		defer stateFile.Close()
	} else {
		log.Fatalf("File already exists - won't truncate. TODO maybe we should try to replay it??")
	}
	w := bufio.NewWriter(stateFile)
	comChan := make(chan string)
	go writeHandler(comChan, w)

	// start a websocket-based Real Time API session
	ws, id := slackConnect(os.Args[1])
	log.Println("sked ready, ^C exits")

	// main loop
	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Received message:")
		log.Println(m)

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			parts := strings.Fields(m.Text)
			// command name is first argument
			com_name := parts[1]
			var msg string

			// 'help' is treated specially
			if com_name == "help" {
				msg = helpAction(command_map, parts)
			} else if act, ok := command_map[com_name]; ok {
				// if we know the command...
				// write to command log
				comChan <- m.Text
				c := command{parts[1], parts[2:]}
				msg = act.function(c, &sked_state)
			} else {
				// we don't know the command
				msg = fmt.Sprintln("sorry, that does not compute")
			}
			m.Text = msg
			go postMessage(ws, m)
		}
	}
}

func helpAction(command_map map[string]action, parts []string) string {
	if len(parts) > 2 {
		act, ok := command_map[parts[2]]
		if ok {
			return fmt.Sprintf("```  %v: %v```", parts[2], act.help)
		} else {
			return fmt.Sprintf("Unknown command %v", parts[2])
		}
	} else if len(parts) == 2 {
		help_list := make([]string, len(command_map))
		i := 0
		for com, act := range command_map {
			help_list[i] = fmt.Sprintf("  %8v: %v", com, act.help)
			i += 1
		}
		return "```" + strings.Join(help_list, "\n") + "```"
	}
	return "" // doesn't happen
}

func getCurrent(cc command, s *state) string {
	return "Not currently implemented"
}

func addPerson(cc command, s *state) string {
	name := cc.args[0]
	if _, ok := s.people[name]; ok {
		return "We already have a " + name + " please choose a different name"
	}
	s.people[name] = NewPerson(name)
	return name + " added"
}

func addUnavailable(cc command, s *state) string {
	name := cc.args[0]
	p, ok := s.people[name]
	if !ok {
		return fmt.Sprintf("I don't know anyone named %v", name)
	}
	startDate, err := getDate(cc.args[1])
	var endDate time.Time
	if err != nil {
		return fmt.Sprintf("I had trouble understanding the date %v, please use the format [YYYY]MMDD[HH]", cc.args[1])
	}
	if len(cc.args) == 4 {
		endDate, err = getDate(cc.args[3])
		if err != nil {
			return fmt.Sprintf("I had trouble understanding the date %v, please use the format [YYYY]MMDD[HH]", cc.args[1])
		}
	} else {
		switch len(cc.args[1]) {
		case 4, 8:
			endDate = startDate.Add(time.Hour * 24)
		case 6, 10:
			endDate = startDate.Add(time.Hour)
		}
	}
	aShift, err := NewShift(startDate, endDate)
	if err != nil {
		return fmt.Sprintf("Your end time:%v is before your start time:%v", endDate, startDate)
	}
	p.AddUnavailable(aShift)
	return fmt.Sprintf("Recorded: %v is unavailable from %v to %v", name, startDate, endDate)
}

// Given a string representing
func getDate(dateStr string) (time.Time, error) {
	var date time.Time
	var err error
	switch len(dateStr) {
	case 4:
		date, err = time.Parse("0102", dateStr)
	case 6:
		date, err = time.Parse("010215", dateStr)
	case 8:
		date, err = time.Parse("20060102", dateStr)
	case 10:
		date, err = time.Parse("2006010215", dateStr)
	}
	if err != nil {
		return time.Time{}, err
	}
	return date, nil

}

func list(cc command, s *state) (msg string) {
	people_names := make([]string, len(s.people))
	i := 0
	for name, _ := range s.people {
		people_names[i] = name
		i += 1
	}
	msg = strings.Join(people_names, ", ")
	if len(s.people) == 0 {
		msg = fmt.Sprintf("List is empty")
	}
	return msg
}

func removePerson(cc command, s *state) (msg string) {
	name := cc.args[0]
	_, ok := s.people[name]
	if !ok {
		return fmt.Sprintf("Could not find '%v'", name)
	} else {
		delete(s.people, name)
		return fmt.Sprintf("'%v' was removed from the list!", name)
	}
}

func getSchedule(cc command, s *state) (msg string) {
	sched := s.BuildSchedule(time.Now(), time.Now().Add(time.Hour*24*7*10))
	return "```" + sched.String() + "```"
}
