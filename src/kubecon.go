package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"../utils"
)

func runit(fn func(*utils.Test), t *utils.Test) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, e := os.Pipe()
	if e != nil {
		fmt.Fprintf(os.Stderr, "e: %s\n", e)
	}
	outFile, err := os.Create("tmpout")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't output output file: %s\n", err)
	}

	os.Stdout = w
	os.Stderr = w

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Failure: %s\n", r)
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
}

func main() {
	verbose := false
	total := 0
	success := 0

	for _, name := range TestNames {
		fn, ok := TestMap[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing func for name %q\n", name)
		}

		t := utils.NewTest()
		runit(fn, t)

		total = total + 1
		if t.Status() && !verbose {
			fmt.Printf("%s: PASS\n", name)
			success = success + 1
		} else {
			fmt.Printf("  +---- %s ----\n", name)
			fi, _ := os.Stat("tmpout")
			if fi.Size() != 0 {
				bytes, _ := ioutil.ReadFile("tmpout")
				lines := strings.Split(strings.TrimSpace(string(bytes)), "\n")
				for _, l := range lines {
					fmt.Printf("  | %s\n", l)
				}
			}
			fmt.Print("  +----\n")
			fmt.Printf("%s: FAIL\n", name)
		}
		os.Remove("tmpout")
	}

	fmt.Printf("\nResults: %d/%d tests passed\n", success, total)

	if total != success {
		os.Exit(1)
	}
}
