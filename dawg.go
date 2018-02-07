package main

// Directed Acyclic Word Graph
// https://en.wikipedia.org/wiki/Deterministic_acyclic_finite_state_automaton
type DAWG struct {
	Terminal bool
	Edge     map[rune]*DAWG
}

// Add adds s to the graph, creating new nodes and marking
// the terminal as necessary.
func (d *DAWG) Add(s string) {
	for _, r := range s {
		next, ok := d.Edge[r]
		if !ok {
			next = NewDAWG()
			d.Edge[r] = next
		}

		// Assign to the receiver pointer d!
		d = next
	}
	d.Terminal = true
}

// AddRecursive works like Add but uses a recursive implementation.
func (d *DAWG) AddRecursive(s string) {
	if len(s) == 0 {
		d.Terminal = true
		return
	}
	r := []rune(s)[0]
	next, ok := d.Edge[r]
	if !ok {
		next = NewDAWG()
		d.Edge[r] = next
	}

	next.AddRecursive(s[1:])
}

// Contains returns true if s reaches a terminal state starting at d.
func (d *DAWG) Contains(s string) bool {
	for _, r := range s {
		next, ok := d.Edge[r]
		if !ok {
			return false
		}
		// Assign to the receiver pointer instead of recursing:
		d = next
	}
	return d.Terminal
}

func NewDAWG() *DAWG {
	totalNodes += 1
	return &DAWG{
		Edge: map[rune]*DAWG{},
	}
}

type Visitor map[*DAWG]bool

func (v Visitor) Traverse(d *DAWG, f func(e rune, d *DAWG)) {
	for r, g := range d.Edge {
		if v[g] {
			continue
		}
		f(r, g)
		v[g] = true
		v.Traverse(g, f)
	}
}
