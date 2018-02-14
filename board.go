package main

const (
	// Empty represents an unplayed square.
	Empty = rune(' ')
)

var (
	TilePoints map[rune]int
)

func init() {
	TilePoints = map[rune]int{}

	p := map[string]int{
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

//// Game state representation structures.
type Row struct {
	Col []rune
}

// Anchors returns the positions of possible anchor squares in the row.
// An anchor is a square that is vacant and has a played character to
// the right of it.
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

func NewBoard() *Board {
	b := &Board{}
	b.Row = make([]*Row, 15, 15)
	for i, r := range b.Row {
		r = &Row{}
		b.Row[i] = r
		r.Col = make([]rune, 15)
	}
	return b
}

type ScoreType int

const (
	None ScoreType = iota
	DL
	TL
	DW
	TW
)

func ScoreAcross(x, y int, word string) int {
	ret := 0
	wordMult := 1
	for i, r := range word {
		s := ScoreAtConditional(x+i, y)
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

// I was going to write a test to make sure the Bits slices were
// the correct length but then remembered that the compiler can
// check that for me for free!
type scoreBits [8]uint

var (
	TWBits = scoreBits{
		1 | (1 << 7),
		0, 0, 0, 0, 0, 0,
		1 << 7,
	}
	DWBits = scoreBits{
		0,
		1 << 6, 1 << 5, 1 << 4, 1 << 3,
		0, 0,
		1, // Center square is double word score.
	}
	TLBits = scoreBits{
		0,
		1 << 2,
		0, 0, 0,
		(1 << 6) | (1 << 2),
		0, 0,
	}
	DLBits = scoreBits{
		1 << 4,
		0,
		1 << 1,
		1 | (1 << 7),
		0, 0,
		(1 << 1) | (1 << 5),
		1 << 4,
	}
)

func (s scoreBits) At(x, y int) bool {
	row := s[y]
	return row&(1<<uint(7-x)) > 0
}

func ScoreAtBits(x, y int) ScoreType {
	// Symmetric adjustments if x or y > 7 to simplify checks below.
	if x > 7 {
		x = 14 - x
	}

	if y > 7 {
		y = 14 - y
	}

	switch {
	case TWBits.At(x, y):
		return TW
	case DWBits.At(x, y):
		return DW
	case TLBits.At(x, y):
		return TL
	case DLBits.At(x, y):
		return DL
	}

	return None
}

// But you know what? The entire bit mask is 8x8, so we
// could just stuff them all into uint64s.
type scoreInt uint64

var (
	TWInt = scoreInt(0x8000000000000081)
	DWInt = scoreInt(0x100000810204000)
	TLInt = scoreInt(0x440000000400)
	DLInt = scoreInt(0x1022000081020010)
)

func (s scoreInt) At(x, y int) bool {
	return s&scoreInt(1<<(uint(y)*8+uint(7-x))) > 0
}

func ScoreAtInt(x, y int) ScoreType {
	// Symmetric adjustments if x or y > 7 to simplify checks below.
	if x > 7 {
		x = 14 - x
	}

	if y > 7 {
		y = 14 - y
	}

	switch {
	case TWInt.At(x, y):
		return TW
	case DWInt.At(x, y):
		return DW
	case TLInt.At(x, y):
		return TL
	case DLInt.At(x, y):
		return DL
	}

	return None
}

// ScoreAtConditional uses a series of conditional checks to determine what a
// particular square's score multiplier is.
func ScoreAtConditional(x, y int) ScoreType {
	// Symmetric adjustments if x or y > 7 to simplify checks below.
	if x > 7 {
		x = 14 - x
	}

	if y > 7 {
		y = 14 - y
	}

	if x == y {
		if x == 0 {
			return TW
		}
		if x > 0 && x < 5 {
			return DW
		}
		if x == 5 {
			return TL
		}
		if x == 6 {
			return DL
		}
		if x == 7 {
			return DW
		}
	}

	symCheck := func(a, b int) bool {
		return x == a && y == b || x == b && y == a
	}

	switch {
	case symCheck(0, 3) || symCheck(2, 6) || symCheck(3, 7):
		return DL
	case symCheck(0, 7):
		return TW
	case symCheck(1, 5):
		return TL
	case symCheck(3, 7):
		return DW
	}
	return None
}
