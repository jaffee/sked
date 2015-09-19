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
)

type state struct {
	engineers []string
}

type command struct {
	action string
	args   []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: sked <slack-bot-token>\n")
		os.Exit(1)
	}

	command_map := map[string]func(command, *websocket.Conn, Message, *state){
		"stock":   doNothing,
		"current": getCurrent,
		"add":     addPerson,
		"list":    list,
	}
	sked_state := state{make([]string, 0)}
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
			if f, ok := command_map[com_name]; ok {
				// if we know the command...
				c := command{parts[1], parts[2:]}
				f(c, ws, m, &sked_state)

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
	m.Text = s.engineers[0]
	go postMessage(ws, m)
}

func addPerson(cc command, ws *websocket.Conn, m Message, s *state) {
	s.engineers = append(s.engineers, cc.args[0])
	m.Text = cc.args[0] + " added"
	go postMessage(ws, m)
}

func list(cc command, ws *websocket.Conn, m Message, s *state) {
	m.Text = strings.Join(s.engineers, ", ")
	if len(s.engineers) == 0 {
		m.Text = fmt.Sprintln("List is empty")
	}
	go postMessage(ws, m)
}
