package client

import "errors"

var (
	ErrFailedToDial         = errors.New("failed to dial")
	ErrFailedToHandleConn   = errors.New("failed to handle connection")
	ErrFailedToSendMsg      = errors.New("failed to send message")
	ErrFailedToReadResponse = errors.New("failed to read response")
	ErrFailedToHandleResp   = errors.New("failed to handle response")
	ErrFailedToUnmarshal    = errors.New("failed to unmarshal")
)
