package smux

import "net"

type Client struct {
	conn Conn
}

func NewClient(address string, coder coder) *Client {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil
	}
	return &Client{
		conn: *NewConn(conn, coder),
	}
}

func (c *Client) RecvMessage() (Message, error) {
	return c.conn.RecvMessage()
}

func (c *Client) SendMessage(msg *Message) error {
	return c.conn.SendMessage(msg)
}
