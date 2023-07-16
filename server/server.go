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

func (r *Server) Listen(address string) error {
	tcpListener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToListen, err)
		return ErrFailedToListen
	}

	for {
		clientConn, err := tcpListener.Accept()
		if err != nil {
			log.Printf("%s: %v", ErrFailedToAccept, err)
			continue
		}

		go r.handleClientConnection(clientConn)
	}
}

func (r *Server) handleClientConnection(clientConn net.Conn) {
	log.Println("New client connection:", clientConn.RemoteAddr())
	defer clientConn.Close()

	connReader := bufio.NewReader(clientConn)
	for {
		clientRequest, err := connReader.ReadString('\n')
		if err != nil {
			log.Printf("%s: %v", ErrFailedToReadConn, err)
			return
		}

		response, err := r.executeClientRequest(clientRequest, clientConn.RemoteAddr().String())
		if err != nil {
			log.Printf("%s: %v", ErrFailedToProcessReq, err)
			return
		}

		if response != nil {
			err = helpers.SendMessage(*response, clientConn)
			if err != nil {
				log.Printf("%s: %v", ErrFailedToSendMsg, err)
			}
		}
	}
}

func (r *Server) executeClientRequest(clientRequest string, clientAddr string) (*message.Message, error) {
	parsedRequest, err := message.Parse([]byte(clientRequest))
	if err != nil {
		log.Printf("%s: %v", ErrFailedToParse, err)
		return nil, err
	}

	switch parsedRequest.Type {
	case message.ChallengeRequest:
		log.Println("Received Challenge Request")
		return r.handleChallengeRequest(*parsedRequest)

	case message.QuoteRequest:
		log.Printf("Client %s requested quote %s\n", clientAddr, parsedRequest.Data)
		return r.handleQuoteRequest(*parsedRequest, clientAddr)

	default:
		return nil, ErrUnknownRequest
	}
}

func (r *Server) handleChallengeRequest(parsedRequest message.Message) (*message.Message, error) {
	var stamp hashcash.Stamp
	stamp.IssueStamp(parsedRequest.Data, 5)
	randNum, err := strconv.Atoi(stamp.Rand)

	if err != nil {
		log.Printf("%s: %v", ErrFailedToDecodeRand, err)
		return nil, ErrFailedToDecodeRand
	}

	err = r.hashcashService.AddIndicator(randNum)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToAddIndicator, err)
		return nil, ErrFailedToAddIndicator
	}

	marshaledStamp, err := json.Marshal(stamp)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToMarshal, err)
		return nil, ErrFailedToMarshal
	}
	respMsg := message.Message{
		Type: message.ChallengeResponse,
		Data: string(marshaledStamp),
	}

	return &respMsg, nil
}

func (r *Server) handleQuoteRequest(parsedRequest message.Message, clientAddr string) (*message.Message, error) {
	var timeStamp hashcash.Stamp
	err := json.Unmarshal([]byte(parsedRequest.Data), &timeStamp)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToUnmarshal, err)
		return nil, ErrFailedToUnmarshal
	}

	randNum, err := strconv.Atoi(timeStamp.Rand)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToDecodeRand, err)
		return nil, ErrFailedToDecodeRand
	}

	_, err = r.hashcashService.GetIndicator(randNum)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToGetRand, err)
		return nil, ErrFailedToGetRand
	}

	if !timeStamp.IsSolved() {
		log.Printf("%s: %s", ErrChallengeUnsolved, "Challenge unsolved.")
		return nil, ErrChallengeUnsolved
	}

	log.Printf("Client %s successfully computed hashcash %s\n", clientAddr, parsedRequest.Data)

	randomQuote := r.quoteService.GetRandomQuote()

	responseMessage := message.Message{
		Type: message.QuoteResponse,
		Data: randomQuote,
	}

	err = r.hashcashService.RemoveIndicator(randNum)
	if err != nil {
		log.Printf("%s: %v", ErrFailedToRemoveIndicator, err)
		return nil, ErrFailedToRemoveIndicator
	}

	return &responseMessage, nil
}
