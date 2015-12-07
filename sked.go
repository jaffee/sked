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
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type command struct {
	action string
	args   []string
}

type action struct {
	function func(command, *State) string
	help     string
}

func writeHandler(logChan chan string, w *bufio.Writer) {
	for {
		s := <-logChan + "\n"
		fmt.Println("Writing", s, []byte(s))
		n, err := w.Write([]byte(s))
		fmt.Println(n, "bytes written")
		if err != nil {
			log.Printf("Problem while writing, err:%v", err)
			panic(err)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: sked <slack-bot-token> [state-file]\n")
		os.Exit(1)
	}
	gob.Register(Interval{})
	gob.Register(Shift{})

	token := os.Args[1]
	// set up state
	// TODO rename all command funcs to <name>Cmd
	commandMap := map[string]action{
		"current":  action{getCurrent, "Tell me who's scheduled right now"},
		"add":      action{addPerson, "Add a new person to be scheduled. add <name> [ordering_num]"},
		"remove":   action{removePerson, "Remove a person from scheduling. remove <name>"},
		"list":     action{list, "List all the possible people that could be scheduled"},
		"unavail":  action{addUnavailable, "unavail <name> <[YYYY]MMDD[HH]> [to [YYYY]MMDD[HH]]"},
		"schedule": action{getSchedule, "Get the schedule which has been previously built. Or build and return it if it hasn't been built."},
		"build":    action{buildSchedule, "(Re)Build the schedule using the people and availabilities given so far"},
		"edit":     action{editScheduleCmd, "edit <name> [YYYY]<MMDD>[HH] to [YYYY]<MMDD>[HH]"},
		"printCal": action{printCal, "Print in Calendar format (experimental)"},
	}
	skedState := NewState(time.Wednesday)
	err := skedState.Populate()
	if err != nil {
		log.Printf("Error populating from %v. err: %v.", skedState.StorageID, err)
	}

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
		// read file line by line
		// play commands
	}
	w := bufio.NewWriter(stateFile)
	logChan := make(chan string)
	go writeHandler(logChan, w)

	run(logChan, token, commandMap, skedState)
}

func run(logChan chan string, token string, command_map map[string]action, skedState State) {
	// start a websocket-based Real Time API session
	ws, id := slackConnect(token)
	log.Println("sked ready, ^C exits")

	// main loop
	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Printf("Wasn't able to receive a message: error: %v, message: %v\n", err, m)
			continue
		}
		log.Println("Received message:")
		log.Println(m)

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			parts := strings.Fields(m.Text)
			// command name is first argument
			var msg string
			var com_name string
			if len(parts) > 1 {
				com_name = parts[1]
			} else {
				com_name = "help"
			}

			// 'help' is treated specially
			if com_name == "help" {
				msg = helpAction(command_map, parts)
			} else if act, ok := command_map[com_name]; ok {
				// if we know the command...
				// write to command log
				logChan <- strings.Join(parts[1:], " ")
				c := command{parts[1], parts[2:]}
				skedState.Lock()
				msg = act.function(c, &skedState)
				err := skedState.Persist()
				skedState.Unlock()
				if err != nil {
					m.Text = fmt.Sprintf("I'm having trouble persisting my state - err: %v", err)
					go postMessage(ws, m)
				}
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

func getCurrent(cc command, s *State) string {
	w, err := s.Schedule.Current()
	if err != nil {
		return err.Error()
	}
	return w.Identifier()
}

func addPerson(cc command, s *State) string {
	name := cc.args[0]
	var ordering int
	if len(cc.args) > 1 {
		ordering64, err := strconv.ParseInt(cc.args[1], 0, 32)
		if err != nil {
			return fmt.Sprintf("Couldn't understand the number you passed in: %v", cc.args[1])
		}
		ordering = int(ordering64)
	} else {
		ordering = 0
	}
	err := s.AddPerson(name, ordering)
	if err == nil {
		return fmt.Sprintf("%v add with ordering %v", name, s.People[name].Ordering())
	} else {
		return fmt.Sprintf("%v", err)
	}
}

func addUnavailable(cc command, s *State) string {
	name := cc.args[0]
	p, ok := s.People[name]
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
	aInterval, err := NewInterval(startDate, endDate)
	if err != nil {
		return fmt.Sprintf("Your end time:%v is before your start time:%v", endDate, startDate)
	}
	p.AddUnavailable(aInterval)
	return fmt.Sprintf("Recorded: %v is unavailable from %v to %v", name, startDate, endDate)
}

// Given a string representing a date in one of several formats, get
// the Time object
func getDate(dateStr string) (time.Time, error) {
	loc := time.Now().Location()
	var date time.Time
	var err error
	switch len(dateStr) {
	case 4, 6:
		dateStr := fmt.Sprintf("%v%v", time.Now().Year(), dateStr)
		date, err = time.ParseInLocation("20060102", dateStr, loc)
	case 8:
		date, err = time.ParseInLocation("20060102", dateStr, loc)
	case 10:
		date, err = time.ParseInLocation("2006010215", dateStr, loc)
	}
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func list(cc command, s *State) (msg string) {
	people_names := make([]string, len(s.People))
	i := 0
	for name, _ := range s.People {
		people_names[i] = name
		i += 1
	}
	msg = strings.Join(people_names, ", ")
	if len(s.People) == 0 {
		msg = fmt.Sprintf("List is empty")
	}
	return msg
}

func removePerson(cc command, s *State) (msg string) {
	name := cc.args[0]
	_, ok := s.People[name]
	if !ok {
		return fmt.Sprintf("Could not find '%v'", name)
	} else {
		delete(s.People, name)
		return fmt.Sprintf("'%v' was removed from the list!", name)
	}
}

func buildSchedule(cc command, s *State) (msg string) {
	sched := s.BuildSchedule(time.Now(), time.Now().Add(time.Hour*24*7*10))
	s.Schedule = sched
	return "```" + sched.String() + "```"
}

func getSchedule(cc command, s *State) (msg string) {
	if s.Schedule != nil && s.Schedule.NumShifts() > 0 {
		return "```" + s.Schedule.String() + "```"
	} else {
		return buildSchedule(cc, s)
	}
}

func editScheduleCmd(cc command, s *State) (msg string) {
	// parse args - call edit Schedule
	name := cc.args[0]
	start, err := getDate(cc.args[1])
	if err != nil {
		return fmt.Sprintf("%v is not a valid date. Error: %v", cc.args[1], err)
	}
	end, err := getDate(cc.args[3])
	if err != nil {
		return fmt.Sprintf("%v is not a valid date. Error: %v", cc.args[3], err)
	}

	person, ok := s.People[name]
	if !ok {
		return "No one named: " + name
	}
	editSchedule(person, start, end, s)
	return "Schedule was edited"
}

func editSchedule(person *Person, start time.Time, end time.Time, s *State) {
	s.Schedule.AddShift(person, start, end)
}

func startScheduling(cc command, s *State) (msg string) {
	go runSchedule(s)
	return "Schedule started"
}

func printCal(cc command, s *State) (msg string) {
	return "```" + s.Schedule.SPrintCalendar() + "```"
}
