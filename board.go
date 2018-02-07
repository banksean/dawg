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
	b.Row = [15]Row{}
	for _, r := range b.Row {
		r.Col = [15]rune{}
	}
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
	// Diagonals
	if x == y {
		if x == 0 || x == 14 {
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

	if x == 0 || x == 14 {
		if y == 3 || y == 11 {
			return DL
		}
	}
	if x == 1 {
		if y == 5 || y == 9 {
			return TL
		}
	}
	return None
}

/*


 */
