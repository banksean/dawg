package main

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParser(t *testing.T) {
	Convey("basic", t, func() {
		// Create a temporary test file
		testData := `>NYMAN ACEINTL 8D AILMENT 70 70
>WAPNICK ADEINOS 7E ANODISE 76 146`
		tmpFile := "test_game.gcg"
		defer func() {
			// Clean up
			os.Remove(tmpFile)
		}()

		err := os.WriteFile(tmpFile, []byte(testData), 0644)
		So(err, ShouldBeNil)

		events := parseFile(tmpFile)
		So(events, ShouldNotBeNil)
		So(len(events), ShouldEqual, 2)
	})
}
