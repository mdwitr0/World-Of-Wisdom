package client

import "errors"

var (
	ErrFailedToDial      = errors.New("failed to dial server")
	ErrFailedToWrite     = errors.New("failed to write to connection")
	ErrFailedToRead      = errors.New("failed to read from connection")
	ErrFailedToUnmarshal = errors.New("failed to unmarshal data")
	ErrFailedToEstablish = errors.New("failed to establish connection")
)
