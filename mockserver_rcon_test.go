package main

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/gorcon/rcon"
)

const (
	MockAddressRCON  = "127.0.0.1:0"
	MockPasswordRCON = "password"

	MockCommandHelpRCON         = "help"
	MockCommandHelpResponseRCON = "lorem ipsum dolor sit amet"
)

// MockServerRCON is a mock Source RCON Protocol server.
type MockServerRCON struct {
	addr        string
	listener    net.Listener
	connections map[net.Conn]struct{}
	wg          sync.WaitGroup
	mu          sync.Mutex
	errors      chan error
	quit        chan bool
}

// NewMockServer returns a running MockServer or nil if an error occurred.
func NewMockServer() (*MockServerRCON, error) {
	listener, err := net.Listen("tcp", MockAddressRCON)
	if err != nil {
		return nil, err
	}

	server := &MockServerRCON{
		listener:    listener,
		connections: make(map[net.Conn]struct{}),
		errors:      make(chan error, 10),
		quit:        make(chan bool),
	}
	server.addr = server.listener.Addr().String()

	server.wg.Add(1)
	go server.serve()

	return server, nil
}

// Close shuts down the MockServer.
func (s *MockServerRCON) Close() error {
	close(s.quit)

	err := s.listener.Close()

	// Waiting for server connections.
	s.wg.Wait()

	// And close remaining connections.
	s.mu.Lock()
	for c := range s.connections {
		// Close connections and add original error if occurred.
		if err2 := c.Close(); err2 != nil {
			if err == nil {
				err = fmt.Errorf("close connenction error: %s", err2)
			} else {
				err = fmt.Errorf("close connenction error: %s. Previous error: %s", err2, err)
			}
		}
	}
	s.mu.Unlock()

	return err
}

// Addr returns IPv4 string MockServer address.
func (s *MockServerRCON) Addr() string {
	return s.addr
}

// serve handles incoming requests until a stop signal is given with Close.
func (s *MockServerRCON) serve() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isRunning() {
				s.reportError(fmt.Errorf("serve error: %s", err))
			}
			return
		}

		s.wg.Add(1)
		go s.handle(conn)
	}
}

// handle handles incoming client conn.
func (s *MockServerRCON) handle(conn net.Conn) {
	s.mu.Lock()
	s.connections[conn] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.closeConnection(conn)
		s.wg.Done()
	}()

	for {
		request := &rcon.Packet{}
		if _, err := request.ReadFrom(conn); err != nil {
			if err == io.EOF {
				return
			}

			s.reportError(fmt.Errorf("handle read request error: %s", err))
			return
		}

		responseType := rcon.SERVERDATA_RESPONSE_VALUE
		responseID := request.ID
		responseBody := ""

		switch request.Type {
		case rcon.SERVERDATA_AUTH:
			responseType = rcon.SERVERDATA_AUTH_RESPONSE
			if request.Body() != MockPasswordRCON {
				// If authentication was failed, the ID must be assigned to -1.
				responseID = -1
				responseBody = string([]byte{0x00})
			}
		case rcon.SERVERDATA_EXECCOMMAND:
			switch request.Body() {
			case MockCommandHelpRCON:
				responseBody = MockCommandHelpResponseRCON
			default:
				responseBody = "unknown command"
			}
		}

		response := rcon.NewPacket(responseType, responseID, responseBody)
		if err := s.write(conn, responseID, response); err != nil {
			s.reportError(fmt.Errorf("handle write response error: %s", err))
			return
		}
	}
}

// isRunning returns true if MockServer is running and false if is not.
func (s *MockServerRCON) isRunning() bool {
	select {
	case <-s.quit:
		return false
	default:
		return true
	}
}

// write writes packets to conn. Replaces packets ids to mirrored id from request.
func (s *MockServerRCON) write(conn net.Conn, id int32, packets ...*rcon.Packet) error {
	for _, packet := range packets {
		packet.ID = id
		_, err := packet.WriteTo(conn)
		if err != nil {
			return err
		}
	}

	return nil
}

// closeConnection closes a client conn and removes it from connections map.
func (s *MockServerRCON) closeConnection(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := conn.Close(); err != nil {
		s.reportError(fmt.Errorf("close conn error: %s", err))
	}
	delete(s.connections, conn)
}

// reportError writes error to errors channel.
func (s *MockServerRCON) reportError(err error) bool {
	if err == nil {
		return false
	}

	select {
	case s.errors <- err:
		return true
	default:
		fmt.Printf("erros channel is locked: %s\n", err)
		return false
	}
}
