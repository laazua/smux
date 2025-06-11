package smux

import "encoding/json"

type coder interface {
	Encode(msg *Message) ([]byte, error)
	Decode(data []byte) (*Message, error)
}

type JsonCode struct{}

func (c *JsonCode) Encode(msg *Message) ([]byte, error) {
	return json.Marshal(msg)
}

func (c *JsonCode) Decode(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
