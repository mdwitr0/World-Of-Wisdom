package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	hashcash2 "main/internal/hashcash"
	"main/internal/quote"
	"net"
)

type Server struct {
	HashCashService hashcash2.IService
	QuoteService    quote.IService
}

func NewServer(hashCashService hashcash2.IService, quoteService quote.IService) *Server {
	return &Server{
		HashCashService: hashCashService,
		QuoteService:    quoteService,
	}
}

func (s *Server) ListenAndServe(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToListen, err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("%s: %s", ErrFailedToAccept, err.Error())
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Receive initial stamp.
	initialStampBytes, err := reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		log.Printf("%s: %s", ErrFailedToRead, err.Error())
		return
	}

	// Send challenge back immediately without verifying.
	_, err = conn.Write(initialStampBytes)
	if err != nil {
		log.Printf("%s: %s", ErrFailedToWrite, err.Error())
		return
	}

	// Receive solved stamp.
	solvedStampBytes, err := reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		log.Printf("%s: %s", ErrFailedToRead, err.Error())
		return
	}

	var solvedStamp hashcash2.Stamp
	err = json.Unmarshal(solvedStampBytes, &solvedStamp)
	if err != nil {
		log.Printf("%s: %s", ErrFailedToUnmarshal, err.Error())
		return
	}

	// Now verify the solved stamp.
	if solvedStamp.IsSolved() {
		log.Printf("valid hashcash received, sending quote")
		s.sendQuote(conn)
	} else {
		log.Printf("%w", ErrInvalidHashcash)
	}
}

func (s *Server) sendQuote(conn net.Conn) {
	foundQuote := s.QuoteService.GetRandomQuote()

	quoteBytes, err := json.Marshal(foundQuote)
	if err != nil {
		log.Printf("%s: %s", ErrFailedToMarshal, err.Error())
		return
	}

	quoteBytes = append(quoteBytes, '\n')

	_, err = conn.Write(quoteBytes)
	if err != nil {
		log.Printf("%s: %s", ErrFailedToWrite, err.Error())
	}
}
