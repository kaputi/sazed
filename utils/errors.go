package utils

import (
	"log"
	"strings"
)

var debug = false

func CheckErr(err error) {
	if err != nil {
		if debug {
			panic(err)
		}
		log.Fatal(err)
	}
}

type ErrList struct {
	errors []error
}

func NewErrList() *ErrList {
	return &ErrList{}
}

func (e *ErrList) Add(err error) {
	if err != nil {
		e.errors = append(e.errors, err)
	}
}

func (e ErrList) Check() {
	if len(e.errors) <= 0 {
		return
	}

	asString := make([]string, len(e.errors))
	for i, err := range e.errors {
		if debug {
			panic(err)
		}
		asString[i] = err.Error()
	}

	log.Fatal(strings.Join(asString, "\n"))
}
