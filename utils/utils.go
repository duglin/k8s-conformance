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

func RunCmd(args ...string) (string, int) {
	// TODO: add timeouts
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), 1
	}

	return string(output), 0
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

// JsonValue will walk a json object (as a string) for a particular value.
// path is of form:
//   prop
//   prop.prop
//   prop[int].prop
func JsonValue(jsonStr string, path string) (string, error) {
	var obj interface{}
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		panic(fmt.Sprintf("Error parsing json(%s):\n%s", err, jsonStr))
	}
	return ObjValue(obj, path)
}

func YamlValue(yamlStr string, path string) (string, error) {
	var obj interface{}
	err := yaml.Unmarshal([]byte(yamlStr), &obj)
	if err != nil {
		panic(fmt.Sprintf("Error parsing yaml(%s):\n%s", err, yamlStr))
	}
	return ObjValue(obj, path)
}

func IsNotFound(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), "Can't find property:")
}

func ObjValue(obj interface{}, path string) (string, error) {
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
				return fmt.Sprintf("%v", obj), nil
			default:
				res, err := json.Marshal(obj)
				if err != nil {
					return "", fmt.Errorf("Error marshalling %q", obj)
				}
				return string(res), nil
			}
		}
		word := words[pos]

		if word == "." {
			obj, ok = obj.(map[string]interface{})
			if !ok {
				return "", fmt.Errorf("Can't convert %q to a map", obj)
			}
		} else if word == "[" {
			pos = pos + 1
			if pos == len(words) {
				return "", fmt.Errorf("Missing index")
			}
			index, err := strconv.Atoi(words[pos])
			if err != nil {
				return "", fmt.Errorf("Can't convert %q to an int: %s",
					words[pos], err)
			}
			obj, ok = obj.([]interface{})
			if !ok {
				return "", fmt.Errorf("Can't convert %q to an array", obj)
			}

			if index < 0 || index >= len(obj.([]interface{})) {
				return "", fmt.Errorf("Index %q is out of range", index)
			}

			obj = obj.([]interface{})[index]
		} else {
			obj, ok = obj.(map[string]interface{})[word]
			if !ok {
				return "", fmt.Errorf("Can't find property: %s", word)
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
