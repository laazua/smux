package main

import (
	"fmt"
	"io"
	"log/slog"

	"smux"
	"smux/auth"
)

func main() {
	// 添加证书
	auth.CaCertFile = "../certs/ssl/ca.crt"
	auth.ServerCrtFile = "../certs/ssl/server.crt"
	auth.ServerKeyFile = "../certs/ssl/server.key"
	auth.AddServerAuthConfig()

	server := smux.NewServer(":8886", &smux.JsonCode{})

	// handler适配器处理客户端消息
	handler := smux.HandlerFunc(MHandler)

	// // 实现smux.Handler接口处理客户端消息
	// handler := &EchoServer{}

	server.SetHandler(handler)

	// 启动服务器
	if err := server.Start(); err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
		return
	}

	// 等待退出信号
	slog.Info("Press Enter to stop the server...")
	fmt.Scanln()

	// 停止服务器
	server.Stop()
}

// 实现smux.Handler接口处理客户端消息
type EchoServer struct{}

func (e *EchoServer) Handle(conn *smux.Conn) error {
	msg, err := conn.RecvMessage()
	if err != nil {
		if err == io.EOF {
			slog.Info("Client disconnected", slog.Any("remote", conn.GetRemoteAddr()))
		} else {
			slog.Error("Read client failure", slog.Any("remoteAddr", conn.GetRemoteAddr()), slog.String("error", err.Error()))
		}
		return err
	}
	slog.Info("Recv client data", slog.Any("content", msg))
	// 这里直接回显了客户端消息
	return conn.SendMessage(&msg)
}

// 使用handler适配器函数处理客户端消息
func MHandler(conn *smux.Conn) error {
	msg, err := conn.RecvMessage()
	if err != nil {
		if err == io.EOF {
			slog.Info("Client disconnected", slog.Any("remote", conn.GetRemoteAddr()))
		} else {
			slog.Error("Read client failure", slog.Any("remoteAddr", conn.GetRemoteAddr()), slog.String("error", err.Error()))
		}
		return err
	}
	slog.Info("Recv client data", slog.Any("content", msg))
	// 自定义回显数据
	msgData := &smux.Message{"id": msg["id"], "status": "OK"}
	return conn.SendMessage(msgData)
}
