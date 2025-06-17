package smux

import (
	"crypto/tls"
	"log/slog"
	"net"

	"smux/auth"
)

type Client struct {
	conn Conn
}

func NewClient(address string, coder coder) *Client {
	var conn net.Conn
	var err error
	if auth.ClientAuthConfig != nil {
		slog.Info("已启用TLS/SSL双向认证")
		conn, err = tls.Dial("tcp", address, auth.ClientAuthConfig)
	} else {
		slog.Info("未启用TLS/SSL双向认证")
		conn, err = net.Dial("tcp", address)
	}
	if err != nil {
		slog.Error("Create client socket failure", slog.String("error", err.Error()))
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
