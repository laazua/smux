package smux

import (
	"errors"
	"fmt"
	"net"
	"sync"
)



type Server struct {
	address     string
	listener    net.Listener
	handler     Handler
	coder       coder
	wg          sync.WaitGroup
	quit        chan struct{}
	connections sync.Map
}

func NewServer(address string, coder coder) *Server {
	return &Server{
		address: address,
		coder:   coder,
		quit:    make(chan struct{}),
	}
}

func (s *Server) SetHandler(handler Handler) {
	s.handler = handler
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("listen error: %v", err)
	}

	fmt.Printf("server started on %s\n", s.address)

	s.wg.Add(1)
	go s.acceptConn()

	return nil
}

// Stop 停止服务器
func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
	fmt.Println("TCP server stopped")
}

// acceptConnections 接受连接
func (s *Server) acceptConn() {
	defer s.wg.Done()

	for {
		select {
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					fmt.Printf("accept error: %v\n", err)
				}
				continue
			}

			connection := NewConn(conn, s.coder)
			s.connections.Store(conn.RemoteAddr().String(), connection)

			s.wg.Add(1)
			go s.handleConn(connection)
		}
	}
}

// handleConnection 处理连接
func (s *Server) handleConn(conn *Conn) {
	remoteAddress := conn.conn.RemoteAddr()
	defer func() {
		conn.conn.Close()
		s.connections.Delete(remoteAddress.String())
		s.wg.Done()
	}()

	fmt.Printf("Client connected: %s\n", remoteAddress)

	for {
		select {
		case <-s.quit:
			return
		default:
			err := s.handler.Handle(conn)
			if err != nil {
				return
			}
		}
	}
}
