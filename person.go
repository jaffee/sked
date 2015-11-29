package main

type Person struct {
	Name           string
	Unavailability []Intervaler
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

func (p *Person) IsAvailable(i Intervaler) bool {
	for _, uInterval := range p.Unavailability {
		if i.Overlaps(uInterval) {
			return false
		}
	}
	return true
}

func (p *Person) AddUnavailable(i Intervaler) {
	p.Unavailability = append(p.Unavailability, i)
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
