package tests

import (
	"fmt"

	. "../utils"
)

// Pod001 will verify that ...
// Conformant implementations MUST ....
func Pod001(t *Test) {
	Helper()
	fmt.Printf("hi from pod001\n")
}

// Pod002 will verify that ...
// Conformant implementations MUST ....
func Pod002(t *Test) {
	fmt.Printf("hi from pod002\n")
	t.Fail("oh no")
}

func Pod003(t *Test) {
}
