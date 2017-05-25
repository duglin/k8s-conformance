package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"../tests"
	. "../utils"
)

func runTest(name string, fn func(*Test)) bool {
	t := NewTest()

	cwd, err := os.Getwd()
	if err != nil {
		ErrStop("Can't get CWD: %s", err)
	}
	defer func() {
		os.Chdir(cwd) // Ignore err?
	}()

	dir, err := ioutil.TempDir("", "kubecon")
	if err != nil {
		ErrStop("Can't create temp dir: %s", err)
	}
	defer func() {
		os.RemoveAll(dir)
	}()

	err = os.Chdir(dir)
	if err != nil {
		ErrStop("Can't cd to temp dir(%q): %s", dir, err)
	}

	// Save old stdout/stderr so we can re-attach them when we're done
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	r, w, e := os.Pipe()
	if e != nil {
		ErrStop("Can't create a pipe: %s", e)
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

	if !t.Status() || Verbose {
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

func ErrStop(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, f+"\n", args...)
	os.Exit(1)
}

func SetupEnv() {
	Logf("Retrieving current context")
	if out, code := Kubectl("config", "current-context"); code != 0 {
		ErrStop("Can't get current context: %s", out)
	}

	Logf("Creating new namespace: %s", tests.TestNS)
	if out, code := Kubectl("create", "namespace", tests.TestNS); code != 0 {
		ErrStop("Can't create the %q namespace: %s", tests.TestNS, out)
	}

	Logf("Set context %q namespace to %q", tests.TestCtxt, tests.TestNS)
	if out, code := Kubectl("config", "set-context", tests.TestCtxt,
		"--namespace", tests.TestNS); code != 0 {
		ErrStop("Can't get namespace: %s", out)
	}
}

func CleanEnv() {
	Logf("Cleaning namespace: %s", tests.TestNS)
	Kubectl("delete", "deployments", "--all", "--force=true", "--now=true",
		"--namespace", tests.TestNS)
	Kubectl("delete", "replicasets", "--all", "--force=true", "--now=true",
		"--namespace", tests.TestNS)
	Kubectl("delete", "pods", "--all", "--force=true", "--now=true",
		"--namespace", tests.TestNS)

	Wait(60*time.Second, func() (bool, error) {
		out, _ := Kubectl("get", "all", "--namespace", tests.TestNS)
		return strings.Contains(out, "No resources found"),
			errors.New("Didn't clean up all the way")
	})
}

func DeleteEnv() {
	CleanEnv()

	Logf("Deleting namespace: %s", tests.TestNS)
	b, err := Wait(60*time.Second, func() (bool, error) {
		_, code := Kubectl("get", "namespace", tests.TestNS)
		if code != 0 {
			return true, nil
		}
		out, code := Kubectl("delete", "namespace", "--force=true",
			tests.TestNS)
		return false, fmt.Errorf("it just won't go away!\n%s", out)
	})

	if !b {
		fmt.Printf("Error deleting namespace %q: %s\n", tests.TestNS, err)
	}

	Logf("Deleting context: %s", tests.TestCtxt)
	Kubectl("config", "delete-context", tests.TestCtxt)
}

func main() {
	flag.BoolVar(&Verbose, "v", false, "Verbose flag")

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
			ErrStop("Bad test pattern %q: %s", pattern, err)
		}
		testPatterns = append(testPatterns, re)
	}

	flag.Args()

	first := true

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

		fn, ok := TestMap[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing func for name %q\n", name)
			continue
		}

		// Don't setup our env until we know for sure we have a test to run
		if first {
			SetupEnv()
		}

		if runTest(name, fn) {
			success = success + 1
		}
		total = total + 1

		CleanEnv()
		first = false
	}

	// Only delete env if we actually set it up
	if !first {
		DeleteEnv()
	}

	if len(testPatterns) > 0 && total == 0 {
		ErrStop("Test name patterns didn't find any matches")
	}

	fmt.Printf("\nResults: %d/%d tests passed\n", success, total)

	if total != success {
		os.Exit(1)
	}
}
