package tests

import (
	"fmt"

	. "../utils"
)

// Pod001 will verify that simple Pod creation works. The platform MUST
// create the specified Pod and queries to retrieve the Pod's metadata MUST
// return the same values that were used when it wad created. The Pod
// MUST eventually end up in the `Running` state, and then be able to be
// deleted. Deleting a Pod MUST remove it from the platform
func Pod001(t *Test) {
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
