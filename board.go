package main

const (
	// Empty represents an unplayed square.
	// rune(0) is the zero rune value.
	Empty = rune(0)

	ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
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
	Legal(string) bool
}

// CrossChecks returns the list of valid runes that may be placed at
// x, y that will not create a word that j rejects.
func (b *Board) CrossChecks(x, y int, j Judge) []rune {
	ret := []rune{}
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
		return []rune(ALPHABET)
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
		if j.Legal(string(w)) {
			ret = append(ret, r)
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
