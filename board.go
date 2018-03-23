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

func (b *Board) String() string {
	ret := ""
	for _, row := range b {
		ret = ret + row.String() + "\n"
	}
	return ret
}

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
		if y >= len(b) {
			panic(fmt.Sprintf("y %d is greater than board len %d", y, len(b)))
		}
		if c+x >= len(b[y]) {
			panic(fmt.Sprintf("x %d + c %d  is greater than board len %d", x, c, len(b[y])))
		}

		b[y][c+x] = r
	}
}

func (b *Board) PlaceDown(x, y int, word string) *Board {
	b = b.Transpose()
	b.PlaceAcross(y, x, word)
	b = b.Transpose()
	return b
}

func (b *Board) ScoreAcross(x, y int, word string) int {
	ret := 0
	wordMult := 0
	newTilesPlayed := 0
	sidePoints := 0
	for i, r := range word {
		// Discount positions that have already been played.
		// They count towards the base score, but multipliers
		// no longer work and we won't check for side points of
		// other words formed vertically since they've already
		// been used in previous plays.
		//fmt.Printf("checking %d, %d: %s\n", x+i, y, string(b[y][x+i]))
		if b[y][x+i] == '*' {
			// Blanks/wildcard tiles don't contribute the score.
			continue
		}
		if b[y][x+i] != Empty {
			ret = ret + TilePoints[r]
			//fmt.Printf("%s was already played\n", string(r))
			continue
		}
		newTilesPlayed += 1
		s := ScrabbleScores.ScoreAt(x+i, y)
		sp := b.SidePoints(x+i, y, r)
		if sp > 0 {
			switch s {
			case TL:
				sp += TilePoints[r] * 3
			case DL:
				sp += TilePoints[r] * 2
			default:
				sp += TilePoints[r]
			}
			switch s {
			case DW:
				sp *= 2
			case TW:
				sp *= 3
			}
		}
		sidePoints += sp

		switch s {
		case DW:
			wordMult += 2
		case TW:
			wordMult += 3
		}
		switch s {
		case TL:
			ret += TilePoints[r] * 3
		case DL:
			ret = ret + TilePoints[r]*2
		default:
			ret = ret + TilePoints[r]
		}
	}
	if wordMult > 0 {
		ret = ret * wordMult
	}

	// Bingo bonus:
	if newTilesPlayed == 7 {
		//fmt.Printf("bingo\n")
		ret += 50
	}

	return ret + sidePoints
}

func (b *Board) SidePoints(x, y int, r rune) int {
	// Check above and below x, y to see if there are tangential words.
	ret := 0
	startY := y
	endY := y

	// stop when startY hits an empty space or 0
	for ; startY > 0; startY-- {
		r := b[startY-1][x]
		if r == Empty {
			break
		}
		//fmt.Printf("adding %d for %s\n", TilePoints[r], string(r))
		ret += TilePoints[r]
	}

	for ; endY < len(b)-1; endY++ {
		r := b[endY+1][x]
		if r == Empty {
			break
		}
		//fmt.Printf("adding %d for %s\n", TilePoints[r], string(r))
		ret += TilePoints[r]
	}

	//fmt.Printf("sp, starting with %s: %d\n", string(r), ret)
	return ret
}

func (b *Board) ScoreDown(x, y int, word string) int {
	b = b.Transpose()
	return b.ScoreAcross(y, x, word)
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

	// stop when startY hits an empty space or 0
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

func (r Row) String() string {
	ret := ""
	for _, t := range r {
		if t == Empty {
			ret = ret + " "
		} else {
			ret = ret + string(t)
		}
	}
	return ret
}

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
func (b Board) LeftPart(x, y int, partialWord string, node *DAWG, limit int, ra Rack, plays chan Play) {
	b.ExtendRight(x, y, partialWord, node, ra, plays)
	if limit > 0 {
		for r, nextNode := range node.Edge {
			if ra[r] > 0 {
				ra.Remove(r)
				b.LeftPart(x, y, partialWord+string(r), nextNode, limit-1, ra, plays)
				ra.Add(r)
			}
		}
	}
}

func (b Board) ExtendRight(x, y int, partialWord string, node *DAWG, ra Rack, plays chan Play) {
	if b[y][x] == Empty {
		if node.Terminal {
			// Send this on a channel?
			LegalWord(partialWord)
			plays <- Play{x, y, partialWord}
		}
		crossChecks := b.CrossChecks(x, y, rootDAWG)
		for r, nextNode := range node.Edge {
			if ra[r] > 0 && crossChecks[r] {
				ra.Remove(r)
				b.ExtendRight(x+1, y, partialWord+string(r), nextNode, ra, plays)
				ra.Add(r)
			}
		}
	} else {
		l := b[y][x]
		if node.Edge[l] != nil {
			nextNode := node.Edge[l]
			b.ExtendRight(x+1, y, partialWord+string(l), nextNode, ra, plays)
		}
	}
}

func LegalWord(s string) {
	fmt.Printf("legal word: %q\n", s)
}

func (b Board) GenerateRowMoves(y int, ra Rack, rootNode *DAWG) chan Play {
	ret := make(chan Play)
	row := b[y]
	anchors := row.Anchors()
	go func() {
		for _, x := range anchors {
			limit := row.LeftMax(x)
			b.LeftPart(x, y, "", rootNode, limit, ra, ret)
		}
		close(ret)
	}()
	return ret
}

func (r Row) LeftMax(x int) int {
	ret := 0
	for i := x; r[i] == Empty; i-- {
		ret++
	}

	return ret
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
