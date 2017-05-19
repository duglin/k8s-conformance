package tests

import (
	"fmt"

	. "../utils"
)

// ReplicaSet001 will ...
func ReplicaSet001(t *Test) {
	fmt.Printf("hi from rs 001\n")
	t.Fail("bad bad bad")
}

// ReplicaSet002 will verify that ...
// Conformant implementations MUST ....
func ReplicaSet002(t *Test) {
	fmt.Printf("hi from rs002\n")
	t.Log("just logging")
	t.Fail("because")
}
