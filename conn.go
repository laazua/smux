package smux

import (
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"
)

type Conn struct {
	writeLock sync.Mutex
	conn      net.Conn
	coder     coder
	sig       Signature
}

func NewConn(conn net.Conn, coder coder) *Conn {
	return &Conn{
		conn:  conn,
		coder: coder,
		sig:   Signature{},
	}
}

func (c *Conn) RecvMessage() (Message, error) {
	c.conn.SetReadDeadline(time.Now().Add(MsgReadTimeout))

	lenBuffer := make([]byte, MsgHeadLen)
	if _, err := io.ReadFull(c.conn, lenBuffer); err != nil {
		return nil, err
	}
	bodyLength := binary.BigEndian.Uint32(lenBuffer)

	if bodyLength > MaxBodySize {
		return nil, errors.New("数据包过大")
	}

	bodyBuffer := make([]byte, bodyLength)
	if _, err := io.ReadFull(c.conn, bodyBuffer); err != nil {
		return nil, err
	}

	decrypted, err := c.sig.Decrypt(bodyBuffer)
	if err != nil {
		return nil, err
	}

	message, err := c.coder.Decode(decrypted)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Conn) SendMessage(msg *Message) error {
	c.conn.SetWriteDeadline(time.Now().Add(MsgWriteTimeout))
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	// 编码并加密
	bodyData, err := c.coder.Encode(msg)
	if err != nil {
		return err
	}
	encryptData, err := c.sig.Encrypt(bodyData)
	if err != nil {
		return err
	}

	// 正确写入加密数据长度
	lenBuffer := make([]byte, MsgHeadLen)
	binary.BigEndian.PutUint32(lenBuffer, uint32(len(encryptData)))

	if _, err := c.conn.Write(lenBuffer); err != nil {
		return err
	}
	if _, err := c.conn.Write(encryptData); err != nil {
		return err
	}

	slog.Info("Socket send data success")
	return nil
}

func (c *Conn) GetRemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
