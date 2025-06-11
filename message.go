package smux

import "time"

const (
	MsgHeadLen = 4

	MsgReadTimeout  = 60 * time.Second
	MsgWriteTimeout = 60 * time.Second
)

type Message struct {
	Id   uint64
	Body []byte
}
