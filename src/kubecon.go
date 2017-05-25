package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"../utils"
)

var verbose = false

func runTest(name string, fn func(*utils.Test)) bool {
	t := utils.NewTest()

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get CWD: %s\n", err)
		os.Exit(1)
	}
	defer func() {
		os.Chdir(cwd) // Ignore err?
	}()

	dir, err := ioutil.TempDir("", "kubecon")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't create temp dir: %s\n", err)
		os.Exit(1)
	}
	defer func() {
		os.RemoveAll(dir)
	}()

	err = os.Chdir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't cd to temp dir(%q): %s\n", dir, err)
		os.Exit(1)
	}

	// Save old stdout/stderr so we can re-attach them when we're done
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	r, w, e := os.Pipe()
	if e != nil {
		fmt.Fprintf(os.Stderr, "Can't create a pipe: %s\n", e)
		os.Exit(1)
	}

	outFile, err := os.Create("stdout") // Should be in temp dir, auto-removed
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't output output file(%q): %s\n",
			"stdout", err)
	}

	// fmt.Printf("Running %q in %q\n", name, dir)

	os.Stdout = w
	os.Stderr = w

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Failure:\n%s\n", r)
				t.SetMessage(fmt.Sprintf("%s", r))
				t.SetStatus(false)
			}
			w.Close()
		}()
		fn(t)
	}()

	io.Copy(outFile, r)

	os.Stdout = oldStdout
	os.Stderr = oldStderr

	w.Close()
	r.Close()
	outFile.Close()

	if t.Status() {
		fmt.Printf("PASS: %s\n", name)
	} else {
		fmt.Printf("FAIL: %s\n", name)
	}

	if !t.Status() || verbose {
		fi, _ := os.Stat("stdout")
		indent := ""
		if fi != nil && fi.Size() != 0 {
			bytes, _ := ioutil.ReadFile("stdout")
			lines := strings.Split(strings.TrimSpace(string(bytes)), "\n")
			first := "|"
			chop := 80 - len(indent) - 2
			for _, l := range lines {
				for len(l) > chop {
					fmt.Printf("%s%s %s", indent, first, l[0:chop])
					l = l[chop:]
					first = "|"
				}
				fmt.Printf("%s%s %s\n", indent, first, l)
				first = "|"
			}
			fmt.Printf("%s+-------------------\n", indent)
		}
	}

	return t.Status()
}

func main() {
	flag.BoolVar(&verbose, "v", false, "Verbose flag")

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

	testPatterns := []*regexp.Regexp{}
	for _, pattern := range flag.Args() {
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Bad test pattern %q: %s\n", pattern, err)
			os.Exit(1)
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

		if fn, ok := TestMap[name]; !ok {
			fmt.Fprintf(os.Stderr, "Missing func for name %q\n", name)
		} else if runTest(name, fn) {
			success = success + 1
		}
		total = total + 1
	}

	if len(testPatterns) > 0 && total == 0 {
		fmt.Fprintln(os.Stderr, "Test name patterns didn't find any matches")
		os.Exit(1)
	}

	fmt.Printf("\nResults: %d/%d tests passed\n", success, total)

	if total != success {
		os.Exit(1)
	}
}
