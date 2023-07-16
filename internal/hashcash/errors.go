package hashcash

import "errors"

var (
	ErrIndicatorNotFound       = errors.New("indicator not found")
	ErrorMaxIterationsExceeded = errors.New("max iterations exceeded")
)
