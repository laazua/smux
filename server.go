package smux

import (
	"crypto/tls"
	"errors"
	"log/slog"
	"net"
	"sync"

	"smux/auth"
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
	if auth.ServerAuthConfig != nil {
		slog.Info("已启用TLS/SSL双向认证")
		s.listener, err = tls.Listen("tcp", s.address, auth.ServerAuthConfig)
	} else {
		slog.Info("未启用TLS/SSL双向认证")
		s.listener, err = net.Listen("tcp", s.address)
	}
	if err != nil {
		return err
	}

	slog.Info("Server started", slog.String("listen", s.address))

	s.wg.Add(1)
	go s.acceptConn()

	return nil
}

// Stop 停止服务器
func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
	slog.Info("TCP server stopped")
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
					slog.Error("Accept failure", slog.String("error", err.Error()))
				}
				continue
			}

			tlsConn, ok := conn.(*tls.Conn)
			if ok {
				err := tlsConn.Handshake()
				if err != nil {
					slog.Error("Server handshake failure", slog.String("error", err.Error()))
					return
				}
				state := tlsConn.ConnectionState()
				for _, v := range state.PeerCertificates {
					slog.Info("Client cert info", slog.Any("subject", v.Subject))
				}
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

	slog.Info("Client connected", slog.Any("remoteAddr", remoteAddress))

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
