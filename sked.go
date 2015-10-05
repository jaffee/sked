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
	"fmt"
	"github.com/jaffee/sked/schedule"
	"log"
	"os"
	"strings"
	"time"
)

const (
	Week = time.Hour * 7 * 24
)

type dateSet map[time.Time]bool

type person struct {
	name string
}

type Shift struct {
	schedule.Shift
	p person
}

func (sh Shift) String() string {
	return fmt.Sprintf("%v: %v to %v", sh.p.name, sh.Start(), sh.End())
}

type state struct {
	people      []person
	unavailable map[string]dateSet
	//	commandHist []command // coming soon!
}

type command struct {
	action string
	args   []string
}

type action struct {
	function func(command, *state) string
	help     string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: sked <slack-bot-token>\n")
		os.Exit(1)
	}

	command_map := map[string]action{
		"current":  action{getCurrent, "Tell me who's scheduled right now"},
		"add":      action{addPerson, "Add a new person to be scheduled"},
		"list":     action{list, "List all the possible people that could be scheduled"},
		"unavail":  action{addUnavailable, "unavail <name> <YYYYMMDD>"},
		"schedule": action{getSchedule, "Build the schedule using the people and availabilities given so far"},
	}

	sked_state := state{make([]person, 0), map[string]dateSet{}}

	// start a websocket-based Real Time API session
	ws, id := slackConnect(os.Args[1])
	log.Println("sked ready, ^C exits")

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
			if com_name == "help" {
				msg = helpAction(command_map, parts)
			} else if act, ok := command_map[com_name]; ok {
				// if we know the command...
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
	if len(s.people) > 0 {
		return s.people[0].name
	} else {
		return "No one is currently scheduled"
	}
}

func addPerson(cc command, s *state) string {
	name := cc.args[0]
	for _, p := range s.people {
		if p.name == name {
			return "We already have a " + name + " please choose a different name"
		}
	}
	s.people = append(s.people, person{cc.args[0]})
	return name + " added"
}

func addUnavailable(cc command, s *state) string {
	name := cc.args[0]
	nameExists := false
	for _, p := range s.people {
		if name == p.name {
			nameExists = true
		}
	}
	if !nameExists {
		return fmt.Sprintf("I don't know anyone named %v", name)
	}
	datestr := cc.args[1]
	date, err := time.Parse("20060102", datestr)
	if err != nil {
		return fmt.Sprintf("I had trouble understanding the date %v, please use the format YYYYMMDD", datestr)
	}
	ds, found := s.unavailable[name]
	if !found {
		ds = dateSet{}
		s.unavailable[name] = ds
	}
	ds[date] = true
	return fmt.Sprintf("Recorded: %v is unavailable on %v", name, date)
}

func list(cc command, s *state) (msg string) {
	people_names := make([]string, len(s.people))
	for i, pers := range s.people {
		people_names[i] = pers.name
	}
	msg = strings.Join(people_names, ", ")
	if len(s.people) == 0 {
		msg = fmt.Sprintf("List is empty")
	}
	return msg
}

func removePerson(cc command, s *state) (msg string) {
	for i, p := range s.people {
		if p.name == cc.args[0] {
			s.people = append(s.people[:i], s.people[i+1:]...)
			return fmt.Sprintf("'%v' was removed from the list!", cc.args[0])
		}
	}
	return fmt.Sprintf("Could not find '%v'", cc.args[0])
}

func getSchedule(cc command, s *state) (msg string) {
	sched := buildSchedule(time.Now(), time.Now().Add(Week*10), time.Wednesday, s)
	sched_strings := make([]string, len(sched))
	for i, t := range sched {
		sched_strings[i] = t.String()
	}
	return "```" + strings.Join(sched_strings, "\n") + "```"
}

func buildSchedule(start time.Time, until time.Time, offset time.Weekday, s *state) []Shift {
	shifts := schedule.GetWeeklyShifts(start, until, offset)
	sched := populateSchedule(shifts, s)
	return sched
}

// Given a list of shifts (start and end times) and a list of people
// with their availabilities, assign a person to each shift as fairly
// as possible.
func populateSchedule(shifts []schedule.Shift, s *state) []Shift {
	pshifts := make([]Shift, len(shifts))
	people := make([]person, len(s.people))
	copy(people, s.people)
	for i, sh := range shifts {
		curShift := Shift{}
		curShift.Shift = sh
		for j, p := range people {
			if isAvailable(sh.Start(), sh.End(), s.unavailable[p.name]) {
				curShift.p = p
				for k := j + 1; k < len(people); k++ {
					people[k-1] = people[k]
				}
				people[len(people)-1] = p
				break
			}
		}
		pshifts[i] = curShift
	}
	return pshifts
}

// Return true if none of the dates in `unavailable` are between start
// and end
func isAvailable(start time.Time, end time.Time, unavailable dateSet) bool {
	for t, _ := range unavailable {
		if t.Before(end) && t.After(start) {
			return false
		}
	}
	return true
}
