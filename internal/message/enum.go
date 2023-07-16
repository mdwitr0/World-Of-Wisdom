package message

type MessageType int

const (
	ChallengeRequest MessageType = iota
	ChallengeResponse
	QuoteRequest
	QuoteResponse
)
