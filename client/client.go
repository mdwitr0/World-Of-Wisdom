package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"main/internal/hashcash"
	"math/rand"
	"net"
	"time"
)

const maxIterations = 10000000

type Quote struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

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
		return fmt.Errorf("%w: %s", ErrFailedToDial, err.Error())
	}
	defer conn.Close()
	log.Printf("connected to %s:%s", client.Hostname, client.Port)

	err = client.handleConnection(conn)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToEstablish, err.Error())
	}

	return nil
}

func (client *Client) handleConnection(conn net.Conn) error {
	stamp := client.initialStamp()
	stampBytes, _ := json.Marshal(stamp)

	_, err := conn.Write(append(stampBytes, '\n'))
	if err != nil {
		log.Printf("%s: %s", ErrFailedToWrite, err.Error())

		return ErrFailedToWrite
	}
	log.Printf("initial stamp sent")

	reader := bufio.NewReader(conn)
	stampBytes, err = reader.ReadBytes('\n')
	if err != nil {
		log.Printf("%s: %s", ErrFailedToRead, err.Error())

		return ErrFailedToRead
	}

	err = json.Unmarshal(stampBytes, &stamp)
	if err != nil {
		log.Printf("%s: %s", ErrFailedToUnmarshal, err.Error())

		return ErrFailedToUnmarshal
	}

	solvedStamp, _ := stamp.ComputeHash(maxIterations)
	stampBytes, _ = json.Marshal(solvedStamp)
	log.Printf("stamp solved, sending back")

	_, err = conn.Write(append(stampBytes, '\n'))
	if err != nil {
		log.Printf("%s: %s", ErrFailedToWrite, err.Error())

		return ErrFailedToWrite
	}

	log.Printf("solved stamp sent")

	quoteBytes, err := reader.ReadBytes('\n')
	if err != nil {
		log.Printf("%s: %s", ErrFailedToRead, err.Error())

		return ErrFailedToRead
	}

	var quote string
	if err := json.Unmarshal(quoteBytes, &quote); err != nil {
		log.Printf("%s: %s", ErrFailedToUnmarshal, err.Error())

		return ErrFailedToUnmarshal
	}

	log.Printf("received quote: %++v", quote)

	return nil
}

func (client *Client) initialStamp() hashcash.Stamp {
	return hashcash.Stamp{
		Version:    1,
		ZerosCount: 4,
		Date:       time.Now().Unix(),
		Resource:   client.Resource,
		Rand:       fmt.Sprint(rand.Intn(100)),
		Counter:    0,
	}
}
