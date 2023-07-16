package server

import "errors"

var (
	ErrFailedToListen    = errors.New("failed to listen on network")
	ErrFailedToAccept    = errors.New("failed to accept connection")
	ErrFailedToRead      = errors.New("failed to read from connection")
	ErrFailedToWrite     = errors.New("failed to write to connection")
	ErrFailedToMarshal   = errors.New("failed to marshal data")
	ErrFailedToUnmarshal = errors.New("failed to unmarshal data")
	ErrInvalidHashcash   = errors.New("invalid hashcash received")
)
