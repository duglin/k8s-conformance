package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ghodss/yaml"
)

func RunCmdSh(cmd string) (string, int) {
	return RunCmd("sh", "-c", cmd)
}

func RunCmd(args ...string) (string, int) {
	// TODO: add timeouts
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), 1
	}

	return string(output), 0
}

func KubectlSh(cmd string) (string, int) {
	return RunCmdSh("kubectl " + cmd)
}

func Kubectl(args ...string) (string, int) {
	a := append([]string{"kubectl"}, args...)

	return RunCmd(a...)
}

func CreateFile(fn string, content string) {
	err := ioutil.WriteFile(fn, []byte(content), 0644)
	if err != nil {
		panic(fmt.Sprintf("Error creating file(%s): %s", fn, content))
	}
}

const NotFoundValue = "****____NOT_FOUND____****"

// JsonValue will walk a json object (as a string) for a particular value.
// path is of form:
//   prop
//   prop.prop
//   prop[int].prop
func JsonValue(jsonStr string, path string) string {
	var obj interface{}
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		panic(fmt.Sprintf("Error parsing json(%s):\n%s", err, jsonStr))
	}

	foundIt, val, err := ObjValue(obj, path)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	if foundIt {
		return val
	}

	return NotFoundValue
}

func YamlValue(yamlStr string, path string) string {
	var obj interface{}
	err := yaml.Unmarshal([]byte(yamlStr), &obj)
	if err != nil {
		panic(fmt.Sprintf("Error parsing yaml(%s):\n%s", err, yamlStr))
	}

	foundIt, val, err := ObjValue(obj, path)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	if foundIt {
		return val
	}

	return NotFoundValue
}

func ObjValue(obj interface{}, path string) (bool, string, error) {
	words := []string{}
	path = strings.TrimSpace(path)
	prop := ""
	for pos := 0; pos-1 != len(path); pos++ {
		if pos != len(path) {
			ch := path[pos]
			if ch != '.' && ch != '[' && ch != ']' {
				prop = prop + string(ch)
				continue
			}
		}

		if prop != "" {
			words = append(words, prop)
		}
		if pos != len(path) && path[pos] != ']' {
			words = append(words, string(path[pos]))
		}
		prop = ""
	}

	var ok bool
	for pos := 0; ; pos = pos + 1 {
		if pos == len(words) {
			switch obj.(type) {
			case string, int, float64, bool:
				return true, fmt.Sprintf("%v", obj), nil
			default:
				res, err := json.Marshal(obj)
				if err != nil {
					return false, "", fmt.Errorf("Error marshalling %q", obj)
				}
				return true, string(res), nil
			}
		}
		word := words[pos]

		if word == "." {
			obj, ok = obj.(map[string]interface{})
			if !ok {
				return false, "", fmt.Errorf("Can't convert %q to a map", obj)
			}
		} else if word == "[" {
			pos = pos + 1
			if pos == len(words) {
				return false, "", fmt.Errorf("Missing index")
			}
			index, err := strconv.Atoi(words[pos])
			if err != nil {
				return false, "", fmt.Errorf("Can't convert %q to an int: %s",
					words[pos], err)
			}
			obj, ok = obj.([]interface{})
			if !ok {
				return false, "", fmt.Errorf("Can't convert %q to an array",
					obj)
			}

			if index < 0 || index >= len(obj.([]interface{})) {
				return false, "", nil
			}

			obj = obj.([]interface{})[index]
		} else {
			obj, ok = obj.(map[string]interface{})[word]
			if !ok {
				return false, "", nil
			}
		}
	}
}

func Wait(timeout time.Duration, fn func() (bool, error)) (bool, error) {
	endTime := time.Now().Unix() + int64(timeout.Seconds())
	var err error
	var b bool
	for time.Now().Unix() < endTime {
		b, err = fn()
		if b == true {
			return b, err
		}
		time.Sleep(1 * time.Second)
	}
	return false, err
}

func WaitPod(timeout time.Duration, name, ns, state string) error {
	b, err := Wait(20*time.Second, func() (bool, error) {
		out, code := Kubectl("get", "pod/"+name, "-o", "json",
			"--namespace", ns)
		if state == "Deleted" {
			err := errors.New("Pod still there")
			if code == 0 {
				return false, err
			}
			return strings.Contains(out, "not found"), err
		}
		if code != 0 {
			return false, fmt.Errorf("Error getting pod: %s", out)
		}
		s := JsonValue(out, "status.phase")
		return s == state, nil
	})
	if b == true {
		return nil
	}
	if err == nil {
		err = fmt.Errorf("Timed-out waiting for container %q",
			name)
	} else {
		err = fmt.Errorf("Timed-out waiting for container %q: %s",
			name, err)
	}
	return err
}

var Verbose bool

func Log(s string) {
	if !Verbose {
		return
	}
	fmt.Print("[ " + s + " ]\n")
}

func Logf(f string, args ...interface{}) {
	if !Verbose {
		return
	}
	Log(fmt.Sprintf(f, args...))
}

type NamespaceInfo struct {
	name string
	used bool
}

var namespaces []*NamespaceInfo

var NSLock1 sync.Mutex
var NSLock2 sync.Mutex

func GetNamespace() (string, error) {
	NSLock1.Lock()
	for _, ns := range namespaces {
		if ns.used == false {
			ns.used = true
			NSLock1.Unlock()
			return ns.name, nil
		}
	}
	NSLock1.Unlock()

	NSLock2.Lock()
	newName := fmt.Sprintf("kcnamespace%d", len(namespaces)+1)
	NSLock1.Lock()
	namespaces = append(namespaces, &NamespaceInfo{
		name: newName,
		used: true,
	})
	NSLock1.Unlock()
	NSLock2.Unlock()

	Logf("Creating new namespace: %s", newName)
	if out, code := Kubectl("create", "namespace", newName); code != 0 {
		return "", fmt.Errorf("Can't create the %q namespace: %s", newName, out)
	}
	return newName, nil
}

func FreeNamespace(name string) {
	NSLock1.Lock()
	for _, ns := range namespaces {
		if ns.name == name {
			NSLock1.Unlock()
			Logf("Cleaning namespace: %s", name)
			Kubectl("delete", "deployments", "--all", "--force=true",
				"--now=true", "--namespace", name)
			Kubectl("delete", "replicasets", "--all", "--force=true",
				"--now=true", "--namespace", name)
			Kubectl("delete", "pods", "--all", "--force=true", "--now=true",
				"--namespace", name)

			_, err := Wait(60*time.Second, func() (bool, error) {
				out, _ := Kubectl("get", "all", "--namespace", name)
				if strings.Contains(out, "No resources found") {
					return true, nil
				}
				return false,
					fmt.Errorf("Didn't clean up all the way:\n%s", out)
			})
			if err != nil {
				Logf("Error deleting namespace(%s): %s", name, err)
				return
			}
			ns.used = false
			return
		}
	}
	NSLock1.Unlock()
}

func DeleteNamespaces() {
	names := []string{}

	for _, ns := range namespaces {
		names = append(names, ns.name)
	}

	Logf("Deleting namespaces: %v", names)
	args := append([]string{"delete", "namespace", "--force=true"}, names[:]...)
	Kubectl(args[:]...)

	for _, ns := range namespaces {
		b, err := Wait(60*time.Second, func() (bool, error) {
			_, code := Kubectl("get", "namespace", ns.name)
			if code != 0 {
				return true, nil
			}
			out, code := Kubectl("delete", "namespace", "--force=true", ns.name)
			return false, fmt.Errorf("it just won't go away!\n%s", out)
		})

		if !b {
			fmt.Fprintf(os.Stderr, "Error deleting namespace %q: %s\n",
				ns.name, err)
		}
	}
}
