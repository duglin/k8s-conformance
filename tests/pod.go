package tests

import (
	"fmt"
	// "strings"

	. "../utils"
)

// Pod001 will verify that simple Pod creation works. The platform MUST
// create the specified Pod and queries to retrieve the Pod's metadata MUST
// return the same values that were used when it wad created. The Pod
// MUST eventually end up in the `Running` state, and then be able to be
// deleted. Deleting a Pod MUST remove it from the platform
func Pod001(t *Test) {
	CreateFile("pod.json", `
        { "apiVersion": "v1",
          "kind": "Pod",
          "metadata" : {
            "name": "pod001"
          },
          "spec": {
            "containers": [{
              "name": "pod1",
              "image": "nginx",
              "imagePullPolicy": "IfNotPresent"
            }]
          }
        }`)

	out, code := Kubectl("create", "-f", "pod.json")
	t.Assertf(code == 0, "Creating pod failed(%d): %s", code, out)

	out, code = Kubectl("get", "pods")
	t.Assertf(code == 0, "Getting pod failed(%d): %s", code, out)

	out, code = Kubectl("get", "pod/pod001", "-o", "yaml")
	t.Assertf(code == 0, "Getting pod failed(%d): %s", code, out)

	v, err := YamlValue(out, "metadata.name")
	t.Logf("metadata.name: %v - %s", v, err)

	m, err := YamlValue(out, "metadata")
	t.Logf("metadata: %v - %s", m, err)

	cs, err := YamlValue(out, "spec.containers")
	t.Logf("spec.containers: %v - %s", cs, err)

	c1, err := YamlValue(out, "spec.containers[0]")
	t.Logf("spec.containers[0]: %v - %s", c1, err)

	c2, err := YamlValue(out, "spec.containers[0].name")
	t.Logf("spec.containers[0].name: %v - %s", c2, err)

	out, code = Kubectl("delete", "pod/pod001")
	t.Assertf(code == 0, "Deleting pod failed(%d): %s", code, out)
}

// Pod002 will verify that ...
// Conformant implementations MUST ....
func Pod002(t *Test) {
	fmt.Printf("hi from pod002\n")
	t.Fail("oh no")
}

func Pod003(t *Test) {
}
