package server

import (
	"bufio"
	"encoding/json"
	"log"
	"main/internal/hashcash"
	"main/internal/message"
	"main/internal/quote"
	"main/internal/shared/helpers"
	"net"
	"strconv"
)

type Server struct {
	hashcashService hashcash.IService
	quoteService    quote.IService
}

func NewServer(hashcashService hashcash.IService, quoteService quote.IService) *Server {
	return &Server{
		hashcashService: hashcashService,
		quoteService:    quoteService,
	}
}

func (s *Server) Listen(address string) error {
	tcpListener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	defer tcpListener.Close()

	for {
		clientConn, err := tcpListener.Accept()
		if err != nil {
			continue
		}

		go s.handleClientConnection(clientConn)
	}
}

func (s *Server) handleClientConnection(clientConn net.Conn) {
	defer clientConn.Close()

	connReader := bufio.NewReader(clientConn)

	for {
		clientRequest, err := connReader.ReadString('\n')
		if err != nil {
			return
		}

		response, err := s.executeClientRequest(clientRequest)
		if err != nil {
			return
		}

		if response != nil {
			_ = helpers.SendMessage(*response, clientConn)
		}
	}
}

func (s *Server) executeClientRequest(clientRequest string) (*message.Message, error) {
	parsedRequest, err := message.Parse(clientRequest)
	if err != nil {
		return nil, err
	}

	switch parsedRequest.Type {
	case message.ChallengeRequest:
		return s.handleChallengeRequest(*parsedRequest)
	case message.QuoteRequest:
		return s.handleQuoteRequest(*parsedRequest)
	default:
		return nil, ErrUnknownRequest
	}
}

func (s *Server) handleChallengeRequest(parsedRequest message.Message) (*message.Message, error) {
	var stamp hashcash.Stamp
	stamp.IssueStamp(parsedRequest.Data, 5)
	randNum, err := strconv.Atoi(stamp.Rand)
	if err != nil {
		return nil, ErrFailedToDecodeRand
	}

	log.Printf("Adding stamp %++v", stamp)

	err = s.hashcashService.AddIndicator(randNum)
	if err != nil {
		return nil, ErrFailedToAddIndicator
	}

	marshaledStamp, err := json.Marshal(stamp)
	if err != nil {
		return nil, ErrFailedToMarshal
	}

	respMsg := message.NewMessage(message.ChallengeResponse, string(marshaledStamp))

	return respMsg, nil
}

func (s *Server) handleQuoteRequest(parsedRequest message.Message) (*message.Message, error) {
	var stamp hashcash.Stamp
	err := json.Unmarshal([]byte(parsedRequest.Data), &stamp)
	if err != nil {
		return nil, ErrFailedToUnmarshal
	}

	log.Printf("Received stamp %++v", stamp)

	randNum, err := strconv.Atoi(stamp.Rand)
	if err != nil {
		return nil, ErrFailedToDecodeRand
	}

	_, err = s.hashcashService.GetIndicator(randNum)
	if err != nil {
		return nil, ErrFailedToGetRand
	}

	if !stamp.IsSolved() {
		return nil, ErrChallengeUnsolved
	}

	randomQuote := s.quoteService.GetRandomQuote()

	responseMessage := message.NewMessage(message.QuoteResponse, randomQuote)

	err = s.hashcashService.RemoveIndicator(randNum)
	if err != nil {
		return nil, ErrFailedToRemoveIndicator
	}

	log.Printf("Response message: %++v", responseMessage)

	return responseMessage, nil
}
