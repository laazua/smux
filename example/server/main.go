package main

import (
	"fmt"
	"io"

	"smux"
)

func main() {
	server := smux.NewServer(":8886", &smux.JsonCode{})
	server.SetHandler(&EchoServer{})
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

// 实现smux.Handler接口
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

	return conn.SendMessage(msg)
}
