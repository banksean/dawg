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
		So(ScoreAtConditional(0, 0), ShouldEqual, TW)
		So(ScoreAtConditional(1, 1), ShouldEqual, DW)
		So(ScoreAtConditional(2, 2), ShouldEqual, DW)
		So(ScoreAtConditional(3, 3), ShouldEqual, DW)
		So(ScoreAtConditional(4, 4), ShouldEqual, DW)
		So(ScoreAtConditional(5, 5), ShouldEqual, TL)
		So(ScoreAtConditional(6, 6), ShouldEqual, DL)
		So(ScoreAtConditional(7, 7), ShouldEqual, DW)
		So(ScoreAtConditional(0, 3), ShouldEqual, DL)
		So(ScoreAtConditional(3, 0), ShouldEqual, DL)
	})

	Convey("symmetry", t, func() {
		for x := 0; x < 15; x++ {
			for y := 0; y < 15; y++ {
				Convey(fmt.Sprintf("%d, %d", x, y), func() {
					So(ScoreAtConditional(x, y), ShouldEqual, ScoreAtConditional(y, x))
				})
			}
		}
	})

	Convey("bits and ints", t, func() {
		for x := 0; x < 15; x++ {
			for y := 0; y < 15; y++ {
				Convey(fmt.Sprintf("%d, %d", x, y), func() {
					So(ScoreAtBits(x, y), ShouldEqual, ScoreAtConditional(x, y))
					So(ScoreAtInt(x, y), ShouldEqual, ScoreAtConditional(x, y))
					So(ScoreAtLiteral(x, y), ShouldEqual, ScoreAtConditional(x, y))
				})
			}
		}
	})
}

func TestScoreAcross(t *testing.T) {
	Convey("basic", t, func() {
		So(ScoreAcross(0, 0, "OH"), ShouldEqual, 15)
		So(ScoreAtConditional(3, 7), ShouldEqual, DL)
		So(ScoreAcross(3, 7, "QUANT"), ShouldEqual, 48)
	})
}

func BenchmarkScoreAtLiteral(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for x := 0; x < 15; x++ {
			for y := 0; y < 15; y++ {
				_ = ScoreAtLiteral(x, y)
			}
		}
	}
}

func BenchmarkScoreAtConditional(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for x := 0; x < 15; x++ {
			for y := 0; y < 15; y++ {
				_ = ScoreAtConditional(x, y)
			}
		}
	}
}

func BenchmarkScoreAtBits(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for x := 0; x < 15; x++ {
			for y := 0; y < 15; y++ {
				_ = ScoreAtBits(x, y)
			}
		}
	}
}

func BenchmarkScoreAtInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for x := 0; x < 15; x++ {
			for y := 0; y < 15; y++ {
				_ = ScoreAtInt(x, y)
			}
		}
	}
}
