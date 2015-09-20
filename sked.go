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
	"golang.org/x/net/websocket"
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
	function func(command, *websocket.Conn, Message, *state)
	help     string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: sked <slack-bot-token>\n")
		os.Exit(1)
	}

	command_map := map[string]action{
		"stock":   action{doNothing, "blah"},
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
			if act, ok := command_map[com_name]; ok {
				// if we know the command...
				c := command{parts[1], parts[2:]}
				act.function(c, ws, m, &sked_state)

			} else {
				// we don't know the command
				m.Text = fmt.Sprintln("sorry, that does not compute")
				go postMessage(ws, m)
			}
		}
	}
}

func doNothing(cc command, ws *websocket.Conn, m Message, s *state) {

}

func getCurrent(cc command, ws *websocket.Conn, m Message, s *state) {
	m.Text = s.people[0]
	go postMessage(ws, m)
}

func addPerson(cc command, ws *websocket.Conn, m Message, s *state) {
	name := cc.args[0]
	for _, p := range s.people {
		if p == name {
			m.Text = "We already have a " + name + " please choose a different name"
			go postMessage(ws, m)
			return
		}
	}
	s.people = append(s.people, cc.args[0])
	m.Text = name + " added"
	go postMessage(ws, m)
}

func addUnavailable(cc command, ws *websocket.Conn, m Message, s *state) {
	name := cc.args[0]
	nameExists := false
	for _, p := range s.people {
		if name == p {
			nameExists = true
		}
	}
	if !nameExists {
		m.Text = fmt.Sprintf("I don't know anyone named %v", name)
		go postMessage(ws, m)
		return
	}
	datestr := cc.args[1]
	date, err := time.Parse("20060102", datestr)
	if err != nil {
		m.Text = fmt.Sprintf("I had trouble understanding the date %v, please use the format YYYYMMDD", datestr)
		go postMessage(ws, m)
		return
	}
	ds, found := s.unavailable[name]
	if !found {
		ds = dateSet{}
		s.unavailable[name] = ds
	}
	ds[date] = true
	m.Text = fmt.Sprintf("Recorded: %v is unavailable on %v", name, date)
	go postMessage(ws, m)
	fmt.Println(s.unavailable)
}

func list(cc command, ws *websocket.Conn, m Message, s *state) {
	m.Text = strings.Join(s.people, ", ")
	if len(s.people) == 0 {
		m.Text = fmt.Sprintln("List is empty")
	}
	go postMessage(ws, m)
}
