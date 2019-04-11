package gocrid

import (
	"errors"
	"strings"
)

type UnallowedOperationError struct {
	opName string
	cond   string
}

var (
	ErrAlreadyLogin error = errors.New("gocrid: user has already login")
	ErrNotLogin     error = errors.New("gocrid: user is NOT login")

	ErrHostBindingDisable error = errors.New("gocrid: host binding is disable")

	ErrRandExit error = errors.New("gocrid: rand is exit")
)

func newUnallowedOperationError(opName string, cond string) error {
	opName = strings.TrimSpace(opName)
	if opName == "" {
		opName = "operation"
	}
	cond = strings.TrimSpace(cond)
	return &UnallowedOperationError{opName: opName, cond: cond}
}

func (uoe *UnallowedOperationError) Error() string {
	var b strings.Builder
	b.Grow(21 + len(uoe.opName) + len(uoe.cond))
	b.WriteString("gocrid: ")
	b.WriteString(uoe.opName)
	b.WriteString("is unallowed")
	if uoe.cond != "" {
		b.WriteRune(' ')
		b.WriteString(uoe.cond)
	}
	return b.String()
}
