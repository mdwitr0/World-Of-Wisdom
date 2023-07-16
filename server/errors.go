package server

import "errors"

var (
	ErrFailedToListen          = errors.New("failed to listen")
	ErrFailedToAccept          = errors.New("failed to accept connection")
	ErrFailedToReadConn        = errors.New("failed to read connection")
	ErrFailedToProcessReq      = errors.New("failed to process request")
	ErrFailedToSendMsg         = errors.New("failed to send message")
	ErrFailedToParse           = errors.New("failed to parse request")
	ErrFailedToAddIndicator    = errors.New("error adding indicator")
	ErrFailedToMarshal         = errors.New("error marshaling timestamp")
	ErrFailedToDecodeRand      = errors.New("error decode rand")
	ErrFailedToGetRand         = errors.New("error get rand from cache")
	ErrChallengeUnsolved       = errors.New("challenge is not solved")
	ErrUnknownRequest          = errors.New("unknown request received")
	ErrFailedToUnmarshal       = errors.New("failed to unmarshal")
	ErrFailedToRemoveIndicator = errors.New("failed to remove indicator")
)
