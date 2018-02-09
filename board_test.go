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
		So(ScoreAt(0, 0), ShouldEqual, TW)
		So(ScoreAt(1, 1), ShouldEqual, DW)
		So(ScoreAt(2, 2), ShouldEqual, DW)
		So(ScoreAt(3, 3), ShouldEqual, DW)
		So(ScoreAt(4, 4), ShouldEqual, DW)
		So(ScoreAt(5, 5), ShouldEqual, TL)
		So(ScoreAt(6, 6), ShouldEqual, DL)
		So(ScoreAt(7, 7), ShouldEqual, DW)

		So(ScoreAt(0, 3), ShouldEqual, DL)
		So(ScoreAt(3, 0), ShouldEqual, DL)
	})

	Convey("symmetry", t, func() {
		for x := 0; x < 15; x++ {
			dlBitsX := uint(0)
			for y := 0; y < 15; y++ {
				if ScoreAt(x, y) == DL {
					dlBitsX = dlBitsX | (1 << uint(14-y))
				}
				Convey(fmt.Sprintf("%d, %d", x, y), func() {
					So(ScoreAt(x, y), ShouldEqual, ScoreAt(y, x))
				})
			}
		}
	})
}

func TestScoreAcross(t *testing.T) {
	Convey("basoc", t, func() {
		So(ScoreAcross(0, 0, "OH"), ShouldEqual, 15)
		So(ScoreAt(3, 7), ShouldEqual, DL)
		So(ScoreAcross(3, 7, "QUANT"), ShouldEqual, 48)
	})

}
