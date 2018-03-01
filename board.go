package main

import (
	"fmt"
)

const (
	// Empty represents an unplayed square.
	// rune(0) is the zero rune value.
	Empty = rune(0)

	ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	TilePoints map[rune]int
	rootDAWG   *DAWG
)

func init() {
	TilePoints = map[rune]int{}

	p := map[string]int{
		// Blank tiles are worth zero points, so yay zero values :)
		"EAIONRTLSU": 1,
		"DG":         2,
		"BCMP":       3,
		"FHVWY":      4,
		"K":          5,
		"JX":         8,
		"QZ":         10,
	}

	for s, n := range p {
		for _, r := range s {
			TilePoints[r] = n
		}
	}
}

// Board is row-major, i.e. [y][x].
type Board [15]Row

// Transpose returns a new Board populated by the
// transposition of b.
func (b *Board) Transpose() *Board {
	a := &Board{}
	for x := range b {
		for y := range b[x] {
			a[x][y] = b[y][x]
		}
	}
	return a
}

func (b *Board) PlaceAcross(x, y int, word string) {
	for c, r := range word {
		// TODO: double check here (or elsewhere)
		// that if y, c+x is already played that it equals r.
		b[y][c+x] = r
	}
}

func (b *Board) ScoreAcross(x, y int, word string) int {
	ret := 0
	wordMult := 1
	// TODO: Discount positions that have already been played!
	for i, r := range word {
		s := ScrabbleScores.ScoreAt(x+i, y)
		switch s {
		case DW:
			wordMult = 2
		case TW:
			wordMult = 3
		}
		switch s {
		case TL:
			ret = ret + TilePoints[r]*3
		case DL:
			ret = ret + TilePoints[r]*2
		default:
			ret = ret + TilePoints[r]
		}
	}
	return ret * wordMult
}

type Judge interface {
	Contains(string) bool
}

// CrossChecks returns the list of valid runes that may be placed at
// x, y that will not create a word that j rejects.
func (b *Board) CrossChecks(x, y int, j Judge) map[rune]bool {
	ret := map[rune]bool{}
	startY := y
	endY := y

	// Stop when startY hits an empty space or 0
	for ; startY > 0; startY-- {
		if b[startY-1][x] == Empty {
			break
		}
	}

	for ; endY < len(b)-1; endY++ {
		if b[endY+1][x] == Empty {
			break
		}
	}

	// If start == end, then this square has empty above and below.
	// So it can be any rune.
	if startY == endY {
		for _, r := range ALPHABET {
			ret[r] = true
		}
		return ret
	}

	// Fill out the test word with the extent around x, y's column.
	// TODO: pre-size this.
	w := []rune{}

	for i := startY; i <= endY; i++ {
		w = append(w, b[i][x])
	}

	// Now for the Judgement!
	for _, r := range ALPHABET {
		w[y-startY] = r
		if j.Contains(string(w)) {
			ret[r] = true
		}
	}

	return ret
}

type Row [15]rune

// Anchors returns the positions of possible anchor squares in the row.
// An anchor is a square that is vacant and has a played character to
// the right of it.
func (r Row) Anchors() []int {
	ret := []int{}
	for i, v := range r {
		if v != Empty {
			continue
		}

		if i < len(r)-1 && r[i+1] != Empty {
			ret = append(ret, i)
		}
	}
	return ret
}

type Play struct {
	x, y int
	word string
}

// More or less literal implementation of pseudocode from the 1988 ACM paper:
func (b Board) LeftPart(x, y int, partialWord string, node *DAWG, limit int, ra Rack) {
	b.ExtendRight(x, y, partialWord, node, ra)
	if limit > 0 {
		for r, nextNode := range node.Edge {
			if ra[r] > 0 {
				ra.Remove(r)
				b.LeftPart(x, y, partialWord+string(r), nextNode, limit-1, ra)
				ra.Add(r)
			}
		}
	}
}

func (b Board) ExtendRight(x, y int, partialWord string, node *DAWG, ra Rack) {
	if b[y][x] == Empty {
		if node.Terminal {
			// Send this on a channel?
			LegalWord(partialWord)
		}
		crossChecks := b.CrossChecks(x, y, rootDAWG)
		for r, nextNode := range node.Edge {
			if ra[r] > 0 && crossChecks[r] {
				ra.Remove(r)
				b.ExtendRight(x+1, y, partialWord+string(r), nextNode, ra)
				ra.Add(r)
			}
		}
	} else {
		l := b[y][x]
		if node.Edge[l] != nil {
			nextNode := node.Edge[l]
			b.ExtendRight(x+1, y, partialWord+string(l), nextNode, ra)
		}
	}
}

func LegalWord(s string) {
	fmt.Printf("legal word: %q\n", s)
}

func (b Board) GenerateMoves(y int, r Rack, d DAWG) []Play {
	ret := []Play{}
	row := b[y]
	anchors := row.Anchors()
	for _, x := range anchors {
		leftParts := row.LeftParts(x, r, d)
		for _, lp := range leftParts {
			rightParts := row.RightParts(x, lp, d)
			for _, rp := range rightParts {
				ret = append(ret, Play{
					x:    x,
					y:    y,
					word: lp + rp,
				})
			}
		}
	}
	return ret
}

func (r Row) LeftParts(x int, ra Rack, d DAWG) []string {
	if x == 0 {
		return nil
	}

	// Check if the squares to the left of x are occupied. If they
	// are, then they form the one and only left part for the
	// anchor at x.
	if r[x-1] != Empty {
		ret := ""
		for i := x - 1; r[i] != Empty; i-- {
			ret = string(r[i]) + ret
		}
		return []string{ret}
	}

	ret := []string{}
	// The squares to the left all have trivial cross checks
	// so we can play any tile from ra in any of them.
	// Just check for all of the prefixes that fit into
	// the free tiles to the left of the anchor that can
	// be constructed from tiles in ra.

	maxLen := 0
	for i := 0; r[i] == Empty; i-- {
		maxLen++
	}

	for i := 1; i < maxLen; i++ {
		// Enumerate all of the prefixes in ra of length i
		// that can be constructed from tiles in ra.
	}

	return ret
}

func (r Row) RightParts(x int, lp string, d DAWG) []string {
	return nil
}

type Rack map[rune]int

func (r Rack) Count() int {
	n := 0
	for _, c := range r {
		n += c
	}
	return n
}

func (r Rack) Add(t rune) {
	if r.Count() > 6 {
		panic("can't add more tiles to rack: " + string(t))
	}

	r[t]++
}

func (r Rack) Remove(t rune) {
	if r[t] <= 0 {
		panic("can't remove tile from rack: " + string(t))
	}

	r[t]--
}

var (
	TileCounts = map[string]int{
		"KJQXZ":       1,
		"BCMPFHVWY":   2,
		"G":           3,
		"DLSU":        4,
		"NRT":         6,
		"O":           8,
		"AI":          9,
		"E":           12,
		string(Empty): 2,
	}
)

type Sack map[rune]int

func NewSack() Sack {
	ret := Sack{}
	for s, c := range TileCounts {
		for _, r := range s {
			ret[r] = c
		}
	}
	return ret
}

func (s Sack) Draw() rune {
	// I *was* bummed that I was going to have to add
	// this package's first import: math/rand. Then
	// I rememberd that the built-in range iterator for
	// maps randomizes the order every time. This isn't
	// as testable as it could be but at least I didn't
	// have to start importing from other packages!
	for r := range s {
		s[r]--
		if s[r] == 0 {
			delete(s, r)
		}
		// Always take the first one returned by
		// the range iterator.
		return r
	}

	panic("tried to draw from an empty Sack")
}

type ScoreType int

const (
	None ScoreType = iota
	DL
	TL
	DW
	TW
)

// Row major, so it's [y][x].
type boardScores [8][8]ScoreType

var (
	ScrabbleScores = boardScores{
		{TW, None, None, DL, None, None, None, TW},
		{None, DW, None, None, None, TL, None, None},
		{None, None, DW, None, None, None, DL, None},
		{DL, None, None, DW, None, None, None, DL},
		{None, None, None, None, DW, None, None, None},
		{None, TL, None, None, None, TL, None, None},
		{None, None, DL, None, None, None, DL, None},
		{TW, None, None, DL, None, None, None, DW},
	}
)

func (b boardScores) ScoreAt(x, y int) ScoreType {
	// Symmetric adjustments if x or y > 7 to simplify checks below.
	if x > 7 {
		x = 14 - x
	}

	if y > 7 {
		y = 14 - y
	}
	return b[y][x]
}
