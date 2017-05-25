package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
)

func RunCMD(args ...string) (string, int) {
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

	return RunCMD(a...)
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

func ObjValue(obj interface{}, path string) (string, error) {
	words := []string{}
	path = strings.TrimSpace(path)
	prop := ""
	for pos := 0; pos-1 != len(path); pos++ {
		// fmt.Printf("pos: %d\n", pos)
		if pos != len(path) {
			ch := path[pos]
			// fmt.Printf("ch: %v\n", string(ch))
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
			str, _ := json.Marshal(obj)
			return string(str), nil
		}
		word := words[pos]

		if word == "." {
			obj = obj.(map[string]interface{})
		} else if word == "[" {
			pos = pos + 1
			index, _ := strconv.Atoi(words[pos])
			obj = obj.([]interface{})[index]
		} else {
			obj, ok = obj.(map[string]interface{})[word]
			if !ok {
				return "", fmt.Errorf("Can't find property: %s", word)
			}
		}
	}
}
