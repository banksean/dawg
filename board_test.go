package main

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewBoard(t *testing.T) {
	Convey("New board is correct size", t, func() {
		b := NewBoard()
		So(len(b.Row), ShouldEqual, 15)
	})
}

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
	Convey("basic", t, func() {
		So(ScoreAcross(0, 0, "OH"), ShouldEqual, 15)
		So(ScrabbleScores.ScoreAt(3, 7), ShouldEqual, DL)
		So(ScoreAcross(3, 7, "QUANT"), ShouldEqual, 48)
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
