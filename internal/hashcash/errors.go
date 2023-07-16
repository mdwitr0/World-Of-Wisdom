package hashcash

import "errors"

var (
	ErrIndicatorNotFound       = errors.New("indicator not found")
	ErrCouldNotAddIndicator    = errors.New("could not add indicator to DB")
	ErrCouldNotGetIndicator    = errors.New("could not get indicator from DB")
	ErrorMaxIterationsExceeded = errors.New("max iterations exceeded")
)
