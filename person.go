package main

import (
	"github.com/jaffee/sked/scheduling"
)

type Person struct {
	Name           string
	Unavailability []scheduling.Shift
	PriorityNum    int
	OrderNum       int
}

func NewPerson(name string) *Person {
	p := Person{
		Name:        name,
		PriorityNum: 0,
	}
	return &p
}

func (p *Person) IsAvailable(s scheduling.Shift) bool {
	for _, pshift := range p.Unavailability {
		if s.Overlaps(pshift) {
			return false
		}
	}
	return true
}

func (p *Person) AddUnavailable(s scheduling.Shift) {
	p.Unavailability = append(p.Unavailability, s)
}

func (p *Person) Identifier() string {
	return p.Name
}

func (p *Person) Priority() int {
	return p.PriorityNum
}

func (p *Person) IncPriority(amnt int) {
	p.PriorityNum += amnt
}

func (p *Person) DecPriority(amnt int) {
	p.PriorityNum -= amnt
}

func (p *Person) Ordering() int {
	return p.OrderNum
}

func (p *Person) SetOrdering(value int) {
	p.OrderNum = value
}

type ByPriority []*Person

func (b ByPriority) Len() int      { return len(b) }
func (b ByPriority) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByPriority) Less(i, j int) bool {
	if b[i].Priority() == b[j].Priority() {
		if b[i].Ordering() == b[j].Ordering() {
			return b[i].Identifier() < b[j].Identifier()
		} else {
			return b[i].Ordering() < b[j].Ordering()
		}
	} else {
		return b[i].Priority() < b[j].Priority()
	}
}
