package tests

import (
	"fmt"
	// "strings"
	"time"

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
            "name": "pod001",
			"namespace": "`+TestNS+`"
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

	out, code = Kubectl("get", "pods", "--namespace", TestNS)
	t.Assertf(code == 0, "Getting pod failed(%d): %s", code, out)

	out, code = Kubectl("get", "po/pod001", "-o", "yaml", "--namespace", TestNS)
	t.Assertf(code == 0, "Getting pod failed(%d): %s", code, out)

	b, err := Wait(20*time.Second, func() (bool, error) {
		out, code = Kubectl("get", "pod/pod001", "-o", "yaml",
			"--namespace", TestNS)
		if code != 0 {
			return false, fmt.Errorf("Error getting pod: %s", out)
		}
		s, err := YamlValue(out, "status.phase")
		return s == "Running", err
	})
	t.Assertf(b, "Timed-out waiting for pod to get ready:\n%s", out)

	n, err := YamlValue(out, "metadata.name")
	t.Assertf(err == nil, "Error parsing yaml: %s\n%s", err, out)
	t.Assertf(n == "pod001", "Wrong pod name(%q), expected %q", n, "pod001")

	i, err := YamlValue(out, "spec.containers[0].image")
	t.Assertf(err == nil, "Error parsing yaml: %s\n%s", err, out)
	t.Assertf(i == "nginx", "Wrong image name(%q), expected %q", i, "nginx")

	out, code = Kubectl("delete", "pod/pod001", "--namespace", TestNS)
	t.Assertf(code == 0, "Deleting pod failed(%d): %s", code, out)
}

// Pod002 will verify that ...
// Conformant implementations MUST ....
func Pod002(t *Test) {
}

func Pod003(t *Test) {
}
