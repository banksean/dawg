package main

import (
	"fmt"
)

const (
	// Empty represents an unplayed square.
	// rune(0) is the zero rune value.
	Empty = rune(0)
)

var (
	TilePoints map[rune]int
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
type Board [15][15]rune

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
			fmt.Printf("anchor square at %d. Square to the right is %q\n", i, r[i+1])
			ret = append(ret, i)
		}
	}
	return ret
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
