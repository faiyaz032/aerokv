package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type Server struct {
	host string
	port string
}

func New(host, port string) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%s", s.host, s.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}
	defer listener.Close()

	fmt.Printf("AeroKV Server running on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

		go s.handleRequest(conn)
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr().String())
			} else {
				fmt.Printf("Error reading from %s: %v\n", conn.RemoteAddr().String(), err)
			}
			return
		}

		fmt.Printf("Received Data: %s", message)
		conn.Write([]byte("Message received\n"))
	}
}
