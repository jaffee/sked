package main

import (
	"testing"
	"time"
)

func TestGetCurrent(t *testing.T) {
	cc := command{}
	s := &state{}
	// TODO implement
}

func TestGetCurrentFail(t *testing.T) {
	cc := command{}
	s := &state{}
	// TODO implement
	curr := getCurrent(cc, s)
	if curr != "No one is currently scheduled" {
		t.Fatalf("getCurrent on empty list failed. Resp: %v", curr)
	}
}

func TestAddPerson(t *testing.T) {
	cc := command{"add", []string{"johann"}}
	s := &state{}
	am := addPerson(cc, s)
	if am != "johann added" {
		t.Fatalf("add Person returned an unexpected message '%v' instead of 'johann added'", am)
	}
	if len(s.people) != 3 || s.people[2].name != "johann" {
		t.Fatalf("johann was not added to the list. List is '%v'", s.people)
	}
}

func TestAddUnavailable(t *testing.T) {
	cc := command{"addUnavailable", []string{"johann", "20150425"}}
	s := &state{}
	// TODO implement, test mutiple dates, test other date types
	uam := addUnavailable(cc, s)
	if uam != "Recorded: johann is unavailable from 2015-04-25 00:00:00 +0000 UTC to 2015-04-26 00:00:00 +0000 UTC" {
		t.Fatalf("Unexpected response from addUnavailable: %v", uam)
	}
}

func TestList(t *testing.T) {
	cc := command{}
	s := &state{}
	msg := list(cc, s)
	if msg != "hello, goodbye" {
		t.Fatalf("getCurrent returned %v instead of %v\n", msg, "hello, goodbye")
	}
	s = &state{}
	msg = list(cc, s)
	if msg != "List is empty" {
		t.Fatalf("getCurrent returned %v instead of %v\n", msg, "List is empty")
	}
}

func TestRemovePerson(t *testing.T) {
	cc := command{"remove", []string{"hello"}}
	s := &state{}
	msg := removePerson(cc, s)
	if len(s.people) != 1 {
		t.Fatalf("Should be 1 person left but there are %v", len(s.people))
	} else if s.people[0].name != "goodbye" {
		t.Fatalf("Person remaining in list should be goodbye, but is %v", s.people[0].name)
	} else if msg != "'hello' was removed from the list!" {
		t.Fatalf("msg is wrong - it is: %v", msg)
	}

	cc = command{"remove", []string{"blah"}}
	s = &state{}
	msg = removePerson(cc, s)
	if len(s.people) != 2 {
		t.Fatalf("Should be 2 people left but there are %v", len(s.people))
	} else if s.people[0].name != "hello" {
		t.Fatalf("First person in list should be hello, but is %v", s.people[0].name)
	} else if msg != "Could not find 'blah'" {
		t.Fatalf("msg is wrong - it is: %v", msg)
	}
}

func TestIsAvailable(t *testing.T) {
	loc, _ := time.LoadLocation("America/Chicago")
	start := time.Date(2015, time.October, 10, 0, 0, 0, 0, loc)
	end := time.Date(2015, time.October, 12, 0, 0, 0, 0, loc)
	// TODO reimplement

}
