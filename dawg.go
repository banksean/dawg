package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

const (
	Empty = rune('-')
)

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

//// Game state representation structures.
type Row struct {
	Col []rune
}

func (r *Row) Anchors() []int {
	ret := []int{}
	for i, v := range r.Col {
		if v != Empty {
			continue
		}

		if i == len(r.Col)-1 {
			ret = append(ret, i)
			continue
		}

		if i < len(r.Col)-1 && r.Col[i+1] != Empty {
			ret = append(ret, i)
		}
	}
	return ret
}

type Board struct {
	Row []*Row
}

func (b *Board) CrossChecks(row int) []map[rune]bool {
	ret := []map[rune]bool{}
	return ret
}

var (
	dictFile   = flag.String("dict", "/usr/share/dict/words", "dictionary file")
	bail       = flag.Int("bail", 0, "bail out after this many lines")
	recurse    = flag.Bool("recurse", false, "use recrsive Add method")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

	totalNodes = 0
)

func init() {
	flag.Parse()
}

func main() {
	byts, err := ioutil.ReadFile(*dictFile)
	if err != nil {
		log.Fatalf("trying to read dict file: %#v", err)
	}

	lines := strings.Split(string(byts), "\n")

	d := NewDAWG()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	for i, line := range lines {
		if *recurse {
			d.AddRecursive(strings.ToLower(line))
		} else {
			d.Add(strings.ToLower(line))
		}

		if i%1000 == 0 {
			fmt.Printf("at line %d\n", i)
		}
		if *bail != 0 && i >= *bail {
			break
		}
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}

	log.Printf("%d total nodes\n", totalNodes)
	//d.Add("do")
	//d.Add("dog")
	fmt.Printf("%+v\n", d.Contains("a"))
	fmt.Printf("%+v\n", d.Contains("d"))
	fmt.Printf("%+v\n", d.Contains("do"))
	fmt.Printf("%+v\n", d.Contains("dog"))
	fmt.Printf("%+v\n", d.Contains("doggo"))
	fmt.Printf("%+v\n", d.Contains("asdasd09u0jasd"))
}
