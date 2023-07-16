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

func NewClient(hostname, port, resource string) *Client {
	return &Client{
		Hostname: hostname,
		Port:     port,
		Resource: resource,
	}
}

func (client *Client) Start() error {
	conn, err := net.Dial("tcp", client.Hostname+":"+client.Port)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToDial, err)
		return ErrFailedToDial
	}
	defer conn.Close()

	err = client.handleConnection(conn)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToHandleConn, err)
		return ErrFailedToHandleConn
	}
	return nil
}

func (client *Client) handleConnection(conn net.Conn) error {
	err := client.requestChallenge(conn)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToSendMsg, err)
		return ErrFailedToSendMsg
	}

	resp, err := receiveResponse(conn)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToReadResponse, err)
		return ErrFailedToReadResponse
	}

	quoteRequest, err := handleChallengeResponse(resp)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToHandleResp, err)
		return ErrFailedToHandleResp
	}

	err = helpers.SendMessage(*quoteRequest, conn)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToSendMsg, err)
		return ErrFailedToSendMsg
	}

	respQuote, err := receiveResponse(conn)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToReadResponse, err)
		return ErrFailedToReadResponse
	}

	quote, err := unmarshallQuote(respQuote)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToUnmarshal, err)
		return ErrFailedToUnmarshal
	}

	log.Printf("Received quote: %s", quote)

	return nil
}

func (client *Client) requestChallenge(conn net.Conn) error {
	msg := message.Message{Type: message.ChallengeRequest, Data: ""}
	return helpers.SendMessage(msg, conn)
}

func receiveResponse(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	resp, err := reader.ReadString('\n')
	return resp, err
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
	err := unmarshallStamp(resp, &stamp)
	if err != nil {
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
	err = json.Unmarshal([]byte(challengeResponseMessage.Data), stamp)
	return err
}

func solveStamp(stamp hashcash.Stamp) (*hashcash.Stamp, error) {
	solved, _ := stamp.ComputeHash(maxIterations)
	return &solved, nil
}

func prepareQuoteRequest(solvedStamp hashcash.Stamp) *message.Message {
	solvedStampMarshalled, _ := json.Marshal(solvedStamp)
	return &message.Message{Type: message.QuoteRequest, Data: string(solvedStampMarshalled)}
}
