package main

import (
	"fmt"
	"os"

	"../utils"
)

func runit(fn func(*utils.Test), s *utils.Test) {
	defer func() {
		if r := recover(); r != nil {
			s.Fail(fmt.Sprintf("%s", r))
		}
	}()
	fn(s)
}

func main() {

	for _, name := range TestNames {
		fn, ok := TestMap[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing func for name %q\n", name)
		}

		t := utils.NewTest()
		runit(fn, t)

		if t.Status() {
			fmt.Printf("%s: PASS\n", name)
		} else {
			fmt.Printf("%s: FAIL\n", name)
		}
	}
}
