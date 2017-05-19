package utils

import (
	"fmt"
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

func (t *Test) Logf(f string, args ...string) {
	fmt.Printf(f+"\n", args)
}

func (t *Test) Fail(msg string) {
	t.status = false
	// t.message = msg
	panic(msg + "\n")
}

func (t *Test) Failf(f string, args ...string) {
	t.status = false
	// t.message = fmt.Sprintf(f, args)
	panic(fmt.Sprintf(f+"\n", args))
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
