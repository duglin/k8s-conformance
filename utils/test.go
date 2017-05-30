package utils

import (
	// "errors"
	"fmt"
	// "io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	// "time"
)

type TestDefinition struct {
	Fn        func(*Test)
	Serialize bool
}

type Test struct {
	Name string
	Fn   func(*Test)
	Dir  string
	NS   string

	Status  bool
	Message string
}

func NewTest(nt Test) *Test {
	t := &Test{
		Name:   "foo",
		Status: true,
	}
	return t
}

func (t *Test) Run() bool {
	var err error

	// Nested
	if t.NS != "" {
		err = os.Chdir(t.Dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't cd to dir(%s): %s", t.Dir, err)
			os.Exit(1)
		}
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "%s\n", r)
				os.Exit(1)
			}
		}()
		t.Fn(t)
		os.Exit(0)
	}

	Logf("Running: %s", t.Name)

	t.Status = true
	t.NS = t.GetNamespace()
	t.Dir, err = ioutil.TempDir("", "kubecon")
	if err != nil {
		t.Status = false
		t.Message = fmt.Sprintf("Can't create temp dir: %s", err)
		return false
	}
	defer func() {
		os.RemoveAll(t.Dir)
	}()

	stdoutFile := t.Dir + "/stdout"

	args := []string{os.Args[0], "-fn=" + t.Name, "-fnNS=" + t.NS,
		"-fnDir=" + t.Dir}
	if Verbose {
		args = append(args, "-v")
	}
	output, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	ioutil.WriteFile(stdoutFile, output, 0600)
	t.Status = (err == nil)

	t.CleanNamespace()

	if t.Status {
		fmt.Printf("PASS: %s\n", t.Name)
	} else {
		fmt.Printf("FAIL: %s\n", t.Name)
	}

	if !t.Status || Verbose {
		fi, _ := os.Stat(stdoutFile)
		indent := ""
		if fi != nil && fi.Size() != 0 {
			bytes, _ := ioutil.ReadFile(stdoutFile)
			lines := strings.Split(strings.TrimSpace(string(bytes)), "\n")
			first := "|"
			chop := 80 - len(indent) - 2
			for _, l := range lines {
				for len(l) > chop {
					fmt.Printf("%s%s %s", indent, first, l[0:chop])
					l = l[chop:]
					first = "|"
				}
				fmt.Printf("%s%s %s\n", indent, first, l)
				first = "|"
			}
			fmt.Printf("%s+-------------------\n", indent)
		}
	}

	return t.Status
}

func (t *Test) Log(msg string) {
	fmt.Print(msg + "\n")
}

func (t *Test) Logf(f string, args ...interface{}) {
	t.Log(fmt.Sprintf(f, args...))
}

func (t *Test) Fail(msg string) {
	if _, file, line, ok := runtime.Caller(1); ok {
		if i := strings.LastIndexAny(file, "/\\"); i >= 0 {
			file = file[i+1:]
		}

		msg = file + ":" + strconv.Itoa(line) + ": " + msg
	}
	t.fail(msg)
}

func (t *Test) Failf(f string, args ...interface{}) {
	if _, file, line, ok := runtime.Caller(1); ok {
		if i := strings.LastIndexAny(file, "/\\"); i >= 0 {
			file = file[i+1:]
		}

		f = file + ":" + strconv.Itoa(line) + ": " + f
	}
	t.fail(fmt.Sprintf(f+"\n", args...))
}

func (t *Test) Assert(b bool, msg interface{}) {
	if b {
		return
	}

	str := fmt.Sprintf("%v", msg)

	if _, file, line, ok := runtime.Caller(1); ok {
		if i := strings.LastIndexAny(file, "/\\"); i >= 0 {
			file = file[i+1:]
		}

		str = file + ":" + strconv.Itoa(line) + ": " + str
	}

	t.fail(str)
}

func (t *Test) Assertf(b bool, f string, args ...interface{}) {
	if b {
		return
	}

	if _, file, line, ok := runtime.Caller(1); ok {
		if i := strings.LastIndexAny(file, "/\\"); i >= 0 {
			file = file[i+1:]
		}

		f = file + ":" + strconv.Itoa(line) + ": " + f
	}

	t.fail(fmt.Sprintf(f, args...))
}

func (t *Test) fail(msg string) {
	t.Status = false
	// t.Message = msg
	panic(msg + "\n")
}

func (t *Test) CreateFile(content string) string {
	f, err := ioutil.TempFile(t.Dir, t.Name)
	if err != nil {
		panic(fmt.Sprintf("Can't create temp file: %s", err))
	}

	if _, err := f.Write([]byte(content)); err != nil {
		os.Remove(f.Name())
		panic(fmt.Sprintf("Error writig to temp file: %s", err))
	}
	return f.Name()
}

func (t *Test) GetNamespace() string {
	ns, err := GetNamespace()
	if err != nil {
		panic(fmt.Sprintf("Can't create a new namespace: %s", err))
	}
	return ns
}

func (t *Test) CleanNamespace() {
	FreeNamespace(t.NS)
}
