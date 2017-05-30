package tests

import (
	// "errors"
	// "strings"
	"time"

	. "../utils"
)

// Pod001 will verify that simple Pod creation works. The platform MUST
// create the specified Pod and queries to retrieve the Pod's metadata MUST
// return the same values that were used when it wad created. The Pod
// MUST eventually end up in the `Running` state, and then be able to be
// deleted. Deleting a Pod MUST remove it from the platform.
// +serialize
func Pod001(t *Test) {
	CreateFile("pod.yaml", `
        { "apiVersion": "v1",
          "kind": "Pod",
          "metadata" : {
            "name": "pod001",
			"namespace": "`+t.NS+`"
          },
          "spec": {
            "containers": [{
              "name": "pod1",
              "image": "nginx",
              "imagePullPolicy": "IfNotPresent"
            }]
          }
        }`)

	out, code := Kubectl("create", "-f", "pod.yaml")
	t.Assertf(code == 0, "Creating pod failed(%d): %s", code, out)

	out, code = Kubectl("get", "pods", "--namespace", t.NS)
	t.Assertf(code == 0, "Getting pod failed(%d): %s", code, out)

	Log(out)
	t.Log(out)

	out, code = KubectlSh("get pod/pod001 -o yaml --namespace " + t.NS)
	t.Assertf(code == 0, "Getting pod failed(%d): %s", code, out)

	err := WaitPod(20*time.Second, "pod001", t.NS, "Running")
	t.Assert(err == nil, err)

	n := YamlValue(out, "metadata.name")
	t.Assertf(n == "pod001", "Wrong pod name(%q), expected %q", n, "pod001")

	i := YamlValue(out, "spec.containers[0].image")
	t.Assertf(i == "nginx", "Wrong image name(%q), expected %q", i, "nginx")

	out, code = KubectlSh("delete pod/pod001 --namespace " + t.NS)
	t.Assertf(code == 0, "Deleting pod failed(%d): %s", code, out)

	err = WaitPod(20*time.Second, "pod001", t.NS, "Deleted")
	t.Assertf(err == nil, "Container still around: %s", err)
}

// Pod002 will verify that ...
// Conformant implementations MUST ....
func Pod002(t *Test) {
}

// Pod003 will do something cool
func Pod003(t *Test) {
}
