package main

const (
	// Empty represents an unplayed square.
	Empty = rune(' ')
)

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

func ScoreAt(x, y int) ScoreType {
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
	}

	symCheck := func(a, b int) bool {
		return x == a && y == b || x == b && y == a
	}

	switch {
	case symCheck(0, 3) || symCheck(2, 6):
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
