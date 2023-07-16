package message

type Type int

const (
	ChallengeRequest Type = iota
	ChallengeResponse
	QuoteRequest
	QuoteResponse
)
