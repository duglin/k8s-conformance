package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
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

func WaitPod(timeout time.Duration, podName string, podNS string) error {
	b, err := Wait(20*time.Second, func() (bool, error) {
		out, code := Kubectl("get", "pod/"+podName, "-o", "json",
			"--namespace", podNS)
		if code != 0 {
			return false, fmt.Errorf("Error getting pod: %s", out)
		}
		s := JsonValue(out, "status.phase")
		return s == "Running", nil
	})
	if b == true {
		return nil
	}
	if err == nil {
		err = fmt.Errorf("Timed-out waiting for container %q to start",
			podName)
	} else {
		err = fmt.Errorf("Timed-out waiting for container %q to start: %s",
			podName, err)
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
