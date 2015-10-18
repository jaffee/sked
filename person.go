package main

import (
	"github.com/jaffee/sked/scheduling"
)

type person struct {
	name           string
	unavailability []scheduling.Shift
	priority       int
}

func NewPerson(name string) *person {
	p := person{
		name:     name,
		priority: 0,
	}
	return &p
}

func (p *person) IsAvailable(s scheduling.Shift) bool {
	for _, pshift := range p.unavailability {
		if s.Overlaps(pshift) {
			return false
		}
	}
	return true
}

func (p *person) AddUnavailable(s scheduling.Shift) {
	p.unavailability = append(p.unavailability, s)
}

func (p *person) Identifier() string {
	return p.name
}

func (p *person) Priority() int {
	return p.priority
}

func (p *person) IncPriority(amnt int) {
	p.priority += amnt
}

func (p *person) DecPriority(amnt int) {
	p.priority -= amnt
}
