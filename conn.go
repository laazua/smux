package smux

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"
)

type Conn struct {
	writeLock sync.Mutex
	conn      net.Conn
	coder     coder
	reader    bufio.Reader
}

func NewConn(conn net.Conn, coder coder) *Conn {
	return &Conn{
		conn:   conn,
		coder:  coder,
		reader: *bufio.NewReader(conn),
	}
}

func (c *Conn) RecvMessage() (Message, error) {
	// 设置超时时间
	c.conn.SetReadDeadline(time.Now().Add(MsgReadTimeout))
	// 读取消息体长度
	headBuffer := make([]byte, MsgHeadLen)
	if _, err := io.ReadFull(&c.reader, headBuffer); err != nil {
		return nil, err
	}
	bodyLength := binary.BigEndian.Uint32(headBuffer)
	bodyBuffer := make([]byte, bodyLength)
	if _, err := io.ReadFull(&c.reader, bodyBuffer); err != nil {
		return nil, err
	}
	// 解码
	message, err := c.coder.Decode(bodyBuffer)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Conn) SendMessage(msg *Message) error {
	// 加锁
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	// 编码
	bodyData, err := c.coder.Encode(msg)
	if err != nil {
		return err
	}
	// 设置超时时间
	c.conn.SetWriteDeadline(time.Now().Add(MsgWriteTimeout))
	// 写入消息长度
	headBuffer := make([]byte, MsgHeadLen)
	binary.BigEndian.PutUint32(headBuffer, uint32(len(bodyData)))
	if _, err := c.conn.Write(headBuffer); err != nil {
		return err
	}
	// 写入消息体内容
	if _, err := c.conn.Write(bodyData); err != nil {
		return err
	}
	return nil
}

func (c *Conn) GetRemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
