package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	. "../utils"
)

func ErrStop(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, f+"\n", args...)
	os.Exit(1)
}

var running = 0
var runningMutex = sync.Mutex{}

func AddRunning() {
	runningMutex.Lock()
	running = running + 1
	runningMutex.Unlock()
}

func DoneRunning() {
	runningMutex.Lock()
	running = running - 1
	runningMutex.Unlock()
}

func main() {
	parallel := 1
	fnNS := ""
	fnDir := ""

	fnName := flag.String("fn", "", "")
	flag.StringVar(&fnNS, "fnNS", "", "")
	flag.StringVar(&fnDir, "fnDir", "", "")
	flag.BoolVar(&Verbose, "v", false, "Verbose flag")
	flag.IntVar(&parallel, "p", 1, "Number of parallel tests to run")

	total := 0
	success := 0

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`Usage: `+os.Args[0]+` [TEST] ...

Run the Kubernetes conformance test suite against the currently configured
Kubernetes environment.

Specify 'TEST', a regular expression, to run seletive tests.
`)

		// flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()

	if parallel < 0 {
		parallel = 1
	}

	if fnName != nil && *fnName != "" {
		td, ok := TestMap[*fnName]
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing func for name %q\n", *fnName)
			os.Exit(1)
		}
		t := Test{
			Name: *fnName,
			Fn:   td.Fn,
			Dir:  fnDir,
			NS:   fnNS,
		}
		if t.Run() {
			os.Exit(0)
		}
		os.Exit(1)
	}

	testPatterns := []*regexp.Regexp{}
	for _, pattern := range flag.Args() {
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			ErrStop("Bad test pattern %q: %s", pattern, err)
		}
		testPatterns = append(testPatterns, re)
	}

	flag.Args()

	for _, name := range TestNames {
		if len(testPatterns) != 0 {
			doit := false
			for _, pattern := range testPatterns {
				if pattern.MatchString(name) {
					doit = true
				}
			}
			if !doit {
				continue
			}
		}

		td, ok := TestMap[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing func for name %q\n", name)
			continue
		}

		t := Test{
			Name: name,
			Fn:   td.Fn,
		}

		for (td.Serialize && running > 0) ||
			(!td.Serialize && running >= parallel) {
			time.Sleep(time.Second)
		}

		AddRunning()
		go func() {
			if t.Run() {
				success = success + 1
			}
			DoneRunning()
		}()

		for td.Serialize && running != 0 {
			time.Sleep(time.Second)
		}

		total = total + 1
	}
	for running > 0 {
		time.Sleep(time.Second)
	}

	DeleteNamespaces()

	if len(testPatterns) > 0 && total == 0 {
		ErrStop("Test name patterns didn't find any matches")
	}

	fmt.Printf("\nResults: %d/%d tests passed\n", success, total)

	if total != success {
		os.Exit(1)
	}
}
