package main

import (
	"flag"
	"fmt"

	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

var (
	dictFile   = flag.String("dict", "/usr/share/dict/words", "dictionary file")
	bail       = flag.Int("bail", 0, "bail out after this many lines")
	recurse    = flag.Bool("recurse", false, "use recrsive Add method")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

	totalNodes = 0
)

func init() {
	flag.Parse()
}

func main() {
	byts, err := os.ReadFile(*dictFile)
	if err != nil {
		log.Fatalf("trying to read dict file: %#v", err)
	}

	lines := strings.Split(string(byts), "\n")

	d := NewDAWG()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	for i, line := range lines {
		if *recurse {
			d.AddRecursive(strings.ToLower(line))
		} else {
			d.Add(strings.ToLower(line))
		}

		if i%1000 == 0 {
			fmt.Printf("at line %d\n", i)
		}
		if *bail != 0 && i >= *bail {
			break
		}
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}

	log.Printf("%d total nodes\n", totalNodes)
	//d.Add("do")
	//d.Add("dog")
	fmt.Printf("%+v\n", d.Contains("a"))
	fmt.Printf("%+v\n", d.Contains("d"))
	fmt.Printf("%+v\n", d.Contains("do"))
	fmt.Printf("%+v\n", d.Contains("dog"))
	fmt.Printf("%+v\n", d.Contains("doggo"))
	fmt.Printf("%+v\n", d.Contains("asdasd09u0jasd"))
}
