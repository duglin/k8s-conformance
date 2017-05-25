package utils

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type Test struct {
	status  bool
	message string
}

func NewTest() *Test {
	t := &Test{
		status: true,
	}
	return t
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

func (t *Test) Status() bool {
	return t.status
}

func (t *Test) SetStatus(s bool) {
	t.status = s
}

func (t *Test) Message() string {
	return t.message
}

func (t *Test) SetMessage(str string) {
	t.message = str
}

func (t *Test) Assert(b bool, msg string) {
	if b {
		return
	}

	if _, file, line, ok := runtime.Caller(1); ok {
		if i := strings.LastIndexAny(file, "/\\"); i >= 0 {
			file = file[i+1:]
		}

		msg = file + ":" + strconv.Itoa(line) + ": " + msg
	}

	t.fail(msg)
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
	t.status = false
	// t.message = msg
	panic(msg + "\n")
}
