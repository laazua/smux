package main

import (
	"fmt"
	"io"

	"smux"
)

func main() {
	server := smux.NewServer(":8886", &smux.JsonCode{})

	// handler适配器处理客户端消息
	handler := smux.HandlerFunc(MHandler)

	// // 实现smux.Handler接口处理客户端消息
	// handler := &EchoServer{}

	server.SetHandler(handler)

	// 启动服务器
	if err := server.Start(); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}

	// 等待退出信号
	fmt.Println("Press Enter to stop the server...")
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
			fmt.Printf("Client disconnected: %s\n", conn.GetRemoteAddr())
		} else {
			fmt.Printf("Read error from %s: %v\n", conn.GetRemoteAddr(), err)
		}
		return err
	}
	fmt.Println("recve messge", msg)
	// 这里直接回显了客户端消息
	return conn.SendMessage(&msg)
}

// 使用handler适配器函数处理客户端消息
func MHandler(conn *smux.Conn) error {
	msg, err := conn.RecvMessage()
	if err != nil {
		if err == io.EOF {
			fmt.Printf("Client disconnected: %s\n", conn.GetRemoteAddr())
		} else {
			fmt.Printf("Read error from %s: %v\n", conn.GetRemoteAddr(), err)
		}
		return err
	}
	// 自定义回显数据
	msgData := &smux.Message{"id": msg["id"], "status": "OK"}
	return conn.SendMessage(msgData)
}
