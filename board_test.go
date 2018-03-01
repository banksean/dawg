package main

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestScoreAt(t *testing.T) {
	Convey("spot checks", t, func() {
		So(ScrabbleScores.ScoreAt(0, 0), ShouldEqual, TW)
		So(ScrabbleScores.ScoreAt(1, 1), ShouldEqual, DW)
		So(ScrabbleScores.ScoreAt(2, 2), ShouldEqual, DW)
		So(ScrabbleScores.ScoreAt(3, 3), ShouldEqual, DW)
		So(ScrabbleScores.ScoreAt(4, 4), ShouldEqual, DW)
		So(ScrabbleScores.ScoreAt(5, 5), ShouldEqual, TL)
		So(ScrabbleScores.ScoreAt(6, 6), ShouldEqual, DL)
		So(ScrabbleScores.ScoreAt(7, 7), ShouldEqual, DW)
		So(ScrabbleScores.ScoreAt(0, 3), ShouldEqual, DL)
		So(ScrabbleScores.ScoreAt(3, 0), ShouldEqual, DL)
	})

	Convey("symmetry", t, func() {
		for x := 0; x < 15; x++ {
			for y := 0; y < 15; y++ {
				Convey(fmt.Sprintf("%d, %d", x, y), func() {
					So(ScrabbleScores.ScoreAt(x, y), ShouldEqual, ScrabbleScores.ScoreAt(y, x))
				})
			}
		}
	})
}

func TestScoreAcross(t *testing.T) {
	Convey("spot checks", t, func() {
		var b Board
		So(b.ScoreAcross(0, 0, "OH"), ShouldEqual, 15)
		So(b.ScoreAcross(3, 7, "QUANT"), ShouldEqual, 48)
	})
}

func TestPlaceAcross(t *testing.T) {
	Convey("basic", t, func() {
		var b Board
		b.PlaceAcross(0, 0, "WHEAT")
		So(b[0][0], ShouldEqual, 'W')
		So(b[0][1], ShouldEqual, 'H')
		So(b[0][2], ShouldEqual, 'E')
		So(b[0][3], ShouldEqual, 'A')
		So(b[0][4], ShouldEqual, 'T')

		b.PlaceAcross(7, 4, "WHEAT")
		So(b[4][7], ShouldEqual, 'W')
		So(b[4][8], ShouldEqual, 'H')
		So(b[4][9], ShouldEqual, 'E')
		So(b[4][10], ShouldEqual, 'A')
		So(b[4][11], ShouldEqual, 'T')
	})
}

type testJudge map[string]bool

func (j testJudge) Contains(s string) bool {
	return j[s]
}

func TestCrossChecks(t *testing.T) {

	allLetters := map[rune]bool{}
	for _, r := range ALPHABET {
		allLetters[r] = true
	}

	Convey("empty", t, func() {
		var b Board
		j := testJudge{}
		// If there are no letters on the board, there are no conflicts.
		for y, _ := range b {
			for x, _ := range b[y] {
				So(b.CrossChecks(x, y, j), ShouldResemble, allLetters)
			}
		}
	})

	Convey("some words played", t, func() {
		var b Board
		b[7][7] = 'A'
		j := testJudge{}

		// To the left and right
		So(b.CrossChecks(6, 7, j), ShouldResemble, allLetters)
		So(b.CrossChecks(8, 7, j), ShouldResemble, allLetters)

		// Above and below
		So(b.CrossChecks(7, 6, j), ShouldResemble, map[rune]bool{})
		So(b.CrossChecks(7, 8, j), ShouldResemble, map[rune]bool{})

		// Now add a word to the dict.
		j["AX"] = true
		So(b.CrossChecks(7, 6, j), ShouldResemble, map[rune]bool{})
		So(b.CrossChecks(7, 8, j), ShouldResemble, map[rune]bool{'X': true})
	})
}

func TestAnchors(t *testing.T) {
	Convey("basic", t, func() {
		var r Row
		So(len(r.Anchors()), ShouldEqual, 0)
		r[4] = 'Q'
		So(r.Anchors(), ShouldResemble, []int{3})
		r[5] = 'I'
		So(r.Anchors(), ShouldResemble, []int{3})
		r[7] = 'K'
		So(r.Anchors(), ShouldResemble, []int{3, 6})

	})
}

func TestTranspose(t *testing.T) {
	Convey("basic", t, func() {
		b := &Board{}
		a := b.Transpose()
		So(b, ShouldResemble, a)

		b[0][0] = 'c'
		a = b.Transpose()
		So(b[0][0], ShouldEqual, 'c')
		So(a[0][0], ShouldEqual, 'c')

		b[0][1] = 'd'
		a = b.Transpose()
		So(b[0][1], ShouldEqual, 'd')
		So(a[1][0], ShouldEqual, 'd')
	})
}

func TestSack(t *testing.T) {
	Convey("basic", t, func() {
		s := NewSack()
		So(len(s), ShouldEqual, 27)

		Convey("draw", func() {
			sum := 0
			for i := 0; i < 100; i++ {
				t := s.Draw()
				sum += TilePoints[t]
			}
			So(sum, ShouldEqual, 187)
		})
	})
}

func BenchmarkScrabbleScoresScoreAt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for x := 0; x < 15; x++ {
			for y := 0; y < 15; y++ {
				_ = ScrabbleScores.ScoreAt(x, y)
			}
		}
	}
}
