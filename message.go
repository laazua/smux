package smux

import "time"

const (
	MsgHeadLen  = 4
	MaxBodySize = 10 * 1024 * 1024

	MsgReadTimeout  = 60 * time.Second
	MsgWriteTimeout = 60 * time.Second
)

// type Message struct {
// 	Id   uint64
// 	Body string
// }

type Message map[string]any
