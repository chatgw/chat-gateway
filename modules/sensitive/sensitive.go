package sensitivemod

import (
	"errors"
	"github.com/king133134/sensfilter"
	"log"
	"os"
	"strings"
)

type Checker struct {
	checker *sensfilter.Search
}

func NewChecker() (*Checker, error) {
	checkerType := os.Getenv("SENSITIVE_TYPE")
	log.Println("Checker Type:" + checkerType)
	var search *sensfilter.Search
	switch checkerType {
	case "string":
		sensitiveWord := os.Getenv("SENSITIVE_WORD")
		wordArray := strings.Split(sensitiveWord, ",")
		search = sensfilter.Strings(wordArray)
	default:
		log.Println("Bad checker type")
		return nil, errors.New("Bad checker type:" + checkerType)
	}
	return &Checker{checker: search}, nil
}

func (c *Checker) HasSense(s []byte) bool {
	return c.checker.HasSens(s)
}