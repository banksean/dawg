package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParser(t *testing.T) {
	Convey("basic", t, func() {
		events := parseFile("1993_wsc_f4_wapnick_nyman.gcg.txt")
		So(events, ShouldNotBeNil)
	})
}
