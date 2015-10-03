package main

import (
	"testing"
	"time"
)

func TestGetCurrent(t *testing.T) {
	s := &state{[]person{person{"hello"}, person{"goodbye"}}, map[string]dateSet{}}
	cc := command{}
	curr := getCurrent(cc, s)
	if curr != "hello" {
		t.Fatalf("getCurrent returned %v instead of %v\n", curr, "hello")
	}
}

func TestGetCurrentFail(t *testing.T) {
	s := &state{[]person{}, map[string]dateSet{}}
	cc := command{}
	curr := getCurrent(cc, s)
	if curr != "No one is currently scheduled" {
		t.Fatalf("getCurrent on empty list failed. Resp: %v", curr)
	}
}

func TestAddPerson(t *testing.T) {
	s := &state{[]person{person{"hello"}, person{"goodbye"}}, map[string]dateSet{}}
	cc := command{"add", []string{"johann"}}
	am := addPerson(cc, s)
	if am != "johann added" {
		t.Fatalf("add Person returned an unexpected message '%v' instead of 'johann added'", am)
	}
	if len(s.people) != 3 || s.people[2].name != "johann" {
		t.Fatalf("johann was not added to the list. List is '%v'", s.people)
	}
}

func TestAddUnavailable(t *testing.T) {
	s := &state{[]person{person{"bill"}, person{"johann"}}, map[string]dateSet{}}
	cc := command{"addUnavailable", []string{"johann", "20150425"}}
	uam := addUnavailable(cc, s)
	if uam != "Recorded: johann is unavailable on 2015-04-25 00:00:00 +0000 UTC" {
		t.Fatalf("Unexpected response from addUnavailable: %v", uam)
	}
	for name, dateset := range s.unavailable {
		if name != "johann" {
			t.Fatalf("Unexpected name in unavailable ")
		}
		for date, bool := range dateset {
			if year, month, day := date.Date(); year != 2015 || month != time.April || day != 25 {
				t.Fatalf("Unexpected date added %v", date)
			}
			if bool != true {
				t.Fatalf("bool should be true in dateSet but isn't")
			}
		}
	}
}

func TestList(t *testing.T) {
	s := &state{[]person{person{"hello"}, person{"goodbye"}}, map[string]dateSet{}}
	cc := command{}
	msg := list(cc, s)
	if msg != "hello, goodbye" {
		t.Fatalf("getCurrent returned %v instead of %v\n", msg, "hello, goodbye")
	}
	s = &state{[]person{}, map[string]dateSet{}}
	msg = list(cc, s)
	if msg != "List is empty" {
		t.Fatalf("getCurrent returned %v instead of %v\n", msg, "List is empty")
	}
}

func TestRemovePerson(t *testing.T) {
	cc := command{"remove", []string{"hello"}}
	s := &state{[]person{{"hello"}, {"goodbye"}}, map[string]dateSet{}}
	msg := removePerson(cc, s)
	if len(s.people) != 1 {
		t.Fatalf("Should be 1 person left but there are %v", len(s.people))
	} else if s.people[0].name != "goodbye" {
		t.Fatalf("Person remaining in list should be goodbye, but is %v", s.people[0].name)
	} else if msg != "'hello' was removed from the list!" {
		t.Fatalf("msg is wrong - it is: %v", msg)
	}

	cc = command{"remove", []string{"blah"}}
	s = &state{[]person{{"hello"}, {"goodbye"}}, map[string]dateSet{}}
	msg = removePerson(cc, s)
	if len(s.people) != 2 {
		t.Fatalf("Should be 2 people left but there are %v", len(s.people))
	} else if s.people[0].name != "hello" {
		t.Fatalf("First person in list should be hello, but is %v", s.people[0].name)
	} else if msg != "Could not find 'blah'" {
		t.Fatalf("msg is wrong - it is: %v", msg)
	}
}
