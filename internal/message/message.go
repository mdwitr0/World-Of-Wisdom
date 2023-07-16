package message

import (
	"encoding/json"
)

type Message struct {
	Type Type   `json:"type"`
	Data string `json:"data"`
}

func (m *Message) ToString() (string, error) {
	msgBytes, err := json.Marshal(m)
	return string(msgBytes), err
}

func Parse(messageString string) (*Message, error) {
	var msg Message
	err := json.Unmarshal([]byte(messageString), &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func NewMessage(messageType Type, data string) *Message {
	return &Message{
		Type: messageType,
		Data: data,
	}
}
