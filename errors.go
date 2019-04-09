package gocrid

import "errors"

var (
	ErrAlreadyLogin error = errors.New("gocrid: user has already login")
	ErrNotLogin     error = errors.New("gocrid: user is NOT login")

	ErrHostBindingDisable error = errors.New("gocrid: host binding is disable")

	ErrRandExit error = errors.New("gocrid: rand is exit")
)
