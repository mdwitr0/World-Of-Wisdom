package client

import (
	"bufio"
	"encoding/json"
	"log"
	"main/internal/hashcash"
	"main/internal/message"
	"main/internal/shared/helpers"
	"net"
)

const maxIterations = 10000000

type Client struct {
	Hostname string
	Port     string
	Resource string
}

type Config struct {
	Hostname string
	Port     string
	Resource string
}

func NewClient(config *Config) *Client {
	return &Client{
		Hostname: config.Hostname,
		Port:     config.Port,
		Resource: config.Resource,
	}
}

func (r *Client) Start() error {
	conn, err := net.Dial("tcp", r.Hostname+":"+r.Port)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := r.handleConnection(conn); err != nil {
		return err
	}
	return nil
}

func (r *Client) handleConnection(conn net.Conn) error {
	if err := r.requestChallenge(conn); err != nil {
		return err
	}

	resp, err := receiveResponse(conn)
	if err != nil {
		return err
	}

	return r.handleChallengeResponse(resp, conn)
}

func (r *Client) handleChallengeResponse(resp string, conn net.Conn) error {
	quoteRequest, err := handleChallengeResponse(resp)
	if err != nil {
		return err
	}

	if err := helpers.SendMessage(*quoteRequest, conn); err != nil {
		return err
	}

	respQuote, err := receiveResponse(conn)
	if err != nil {
		return err
	}

	quote, err := unmarshallQuote(respQuote)
	if err != nil {
		return err
	}

	log.Printf("Received quote: %s", quote)

	return nil
}

func (r *Client) requestChallenge(conn net.Conn) error {
	msg := message.NewMessage(message.ChallengeRequest, "")
	return helpers.SendMessage(*msg, conn)
}

func receiveResponse(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	return reader.ReadString('\n')
}

func unmarshallQuote(respQuote string) (string, error) {
	quoteResponseMessage := message.Message{}
	err := json.Unmarshal([]byte(respQuote), &quoteResponseMessage)
	if err != nil {
		return "", err
	}
	return quoteResponseMessage.Data, nil
}

func handleChallengeResponse(resp string) (*message.Message, error) {
	stamp := hashcash.Stamp{}
	if err := unmarshallStamp(resp, &stamp); err != nil {
		return nil, err
	}
	solvedStamp, err := solveStamp(stamp)
	if err != nil {
		return nil, err
	}
	return prepareQuoteRequest(*solvedStamp), nil
}

func unmarshallStamp(resp string, stamp *hashcash.Stamp) error {
	challengeResponseMessage := message.Message{}
	err := json.Unmarshal([]byte(resp), &challengeResponseMessage)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(challengeResponseMessage.Data), stamp)
}

func solveStamp(stamp hashcash.Stamp) (*hashcash.Stamp, error) {
	solved, _ := stamp.ComputeHash(maxIterations)
	return &solved, nil
}

func prepareQuoteRequest(solvedStamp hashcash.Stamp) *message.Message {
	solvedStampMarshalled, _ := json.Marshal(solvedStamp)
	return &message.Message{Type: message.QuoteRequest, Data: string(solvedStampMarshalled)}
}
