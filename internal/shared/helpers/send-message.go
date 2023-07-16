package helpers

import (
	"errors"
	"fmt"
	"log"
	"main/internal/message"
	"net"
)

var ErrFailedToSendMsg = errors.New("failed to send message")

func SendMessage(msg message.Message, clientConn net.Conn) error {
	messageStr, err := msg.ToString()
	if err != nil {
		log.Printf("%s: %v", ErrFailedToSendMsg, err)
		return ErrFailedToSendMsg
	}

	finalMsg := fmt.Sprintf("%s\n", messageStr)
	_, err = clientConn.Write([]byte(finalMsg))
	if err != nil {
		log.Printf("%s: %v", ErrFailedToSendMsg, err)
		return ErrFailedToSendMsg
	}

	return nil
}
