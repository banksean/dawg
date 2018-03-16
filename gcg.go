package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// See the gcg file format description here:
// http://www.poslfit.com/scrabble/gcg/

type event struct {
	player, rack, word           string
	x, y, score, cumulativeScore int
	across, withdrawal           bool
}

func parseFile(f string) []*event {
	ret := []*event{}
	in, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err.Error())
	}

	lines := strings.Split(string(in), "\n")

	for _, line := range lines {
		if !strings.HasPrefix(line, ">") {
			continue
		}
		evt := parseLine(line)
		ret = append(ret, evt)
	}

	return ret
}

var (
	col = map[rune]int{}
)

func init() {
	for i, r := range "ABCDEFGHIJKLMNO" {
		col[r] = i
	}
}

func parseLine(s string) *event {
	parts := strings.Split(s, " ")
	for i, p := range parts {
		parts[i] = strings.Trim(p, " \t")
	}

	event := &event{
		player:     parts[0],
		rack:       parts[1],
		word:       parts[3],
		across:     false,
		withdrawal: false,
	}

	if _, err := fmt.Sscanf(parts[len(parts)-2], "%d", &event.score); err != nil {
		panic(fmt.Sprintf("parsing %q: %v", parts[4], err))
	}

	if _, err := fmt.Sscanf(parts[len(parts)-1], "%d", &event.cumulativeScore); err != nil {
		panic(fmt.Sprintf("parsing %q: %v", parts[4], err))
	}

	if parts[2] == "--" {
		event.withdrawal = true
		return event
	}

	var c rune
	var i int
	var err error
	pos := parts[2]
	if strings.Contains("ABCDEFGHIJKLMNO", string(pos[0])) {
		c = rune(pos[0])
		pos = pos[1:]
		i, err = strconv.Atoi(pos)
		if err != nil {
			panic(fmt.Sprintf("parsing down %q (%q) err = %v, %v", parts[2], pos, err))
		}
	} else {
		c = rune(pos[len(pos)-1])
		pos = pos[:len(pos)-1]
		i, err = strconv.Atoi(pos)
		if err != nil {
			panic(fmt.Sprintf("parsing across: %q (%q) err = %v, %v", parts[2], pos, err))
		}

		event.across = true
	}

	event.x = col[c]
	event.y = i - 1 // Yeah, they start at 1. :/

	return event
}
