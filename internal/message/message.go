package message

import (
	"encoding/json"
)

type Message struct {
	Type MessageType `json:"type"`
	Data string      `json:"data"`
}

func (m *Message) ToString() (string, error) {
	msgBytes, err := json.Marshal(m)
	return string(msgBytes), err
}

func Parse(msgBytes []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(msgBytes, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
