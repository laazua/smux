package main

import (
	"flag"
	"log/slog"

	"smux"
	"smux/auth"
)

var (
	id      uint64
	message string
	address string
)

func main() {
	// 添加证书
	auth.CaCertFile = "../certs/ssl/ca.crt"
	auth.ClientCrtFile = "../certs/ssl/client.crt"
	auth.ClientKeyFile = "../certs/ssl/client.key"
	auth.AddClientAuthConfig()

	flag.Uint64Var(&id, "i", 0, "message id")
	flag.StringVar(&message, "m", "hello world", "send message")
	flag.StringVar(&address, "a", "localhost:8886", "server address")
	flag.Parse()

	client := smux.NewClient(address, &smux.JsonCode{})
	if client == nil {
		slog.Error("创建socket客户端连接失败")
		return
	}
	defer client.Clean()

	// 自定义数据传递结构(json)
	msgData := &smux.Message{"id": id, "message": message}
	err := client.SendMessage(msgData)
	if err != nil {
		slog.Error("Send message failure", slog.String("error", err.Error()))
		return
	}

	resp, err := client.RecvMessage()
	if err != nil {
		slog.Error("Recv message failure", slog.String("error", err.Error()))
		return
	}
	slog.Info("Recv message", slog.Any("content", resp))
}
