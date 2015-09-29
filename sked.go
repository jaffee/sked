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
	"log"
	"os"
	"strings"
	"time"
)

type dateSet map[time.Time]bool

type state struct {
	people      []string
	unavailable map[string]dateSet
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
		"current": action{getCurrent, "Tell me who's scheduled right now"},
		"add":     action{addPerson, "Add a new person to be scheduled"},
		"list":    action{list, "List all the possible people that could be scheduled"},
		"unavail": action{addUnavailable, "unavail <name> <YYYYMMDD>"},
	}

	sked_state := state{make([]string, 0), map[string]dateSet{}}

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
		return s.people[0]
	} else {
		return "No one is currently scheduled"
	}
}

func addPerson(cc command, s *state) string {
	name := cc.args[0]
	for _, p := range s.people {
		if p == name {
			return "We already have a " + name + " please choose a different name"
		}
	}
	s.people = append(s.people, cc.args[0])
	return name + " added"
}

func addUnavailable(cc command, s *state) string {
	name := cc.args[0]
	nameExists := false
	for _, p := range s.people {
		if name == p {
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
	fmt.Println(s.unavailable)
	return fmt.Sprintf("Recorded: %v is unavailable on %v", name, date)
}

func list(cc command, s *state) (msg string) {
	msg = strings.Join(s.people, ", ")
	if len(s.people) == 0 {
		msg = fmt.Sprintf("List is empty")
	}
	return msg
}
