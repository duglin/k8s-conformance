package utils

import (
	"fmt"
)

type Test struct {
	status  bool
	message string
}

func NewTest() *Test {
	s := &Test{
		status: true,
	}
	return s
}

func (s *Test) Fail(msg string) {
	s.status = false
	s.message = msg
}

func (s *Test) Status() bool {
	return s.status
}

func (s *Test) Message() string {
	return s.message
}

func Helper() {
	fmt.Printf("hi from helper\n")
}
