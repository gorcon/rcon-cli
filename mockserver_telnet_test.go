package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/gorcon/telnet"
)

const (
	MockAddressTELNET  = "127.0.0.1:0"
	MockPasswordTELNET = "password"

	MockCommandHelpTELNET         = "help"
	MockCommandHelpResponseTELNET = "lorem ipsum dolor sit amet"
)

// MockServerTELNET is a mock Source TELNET protocol server.
type MockServerTELNET struct {
	addr        string
	listener    net.Listener
	connections map[net.Conn]struct{}
	wg          sync.WaitGroup
	mu          sync.Mutex
	errors      chan error
	quit        chan bool
}

// NewMockServerTELNET returns a running MockServer or nil if an error occurred.
func NewMockServerTELNET() (*MockServerTELNET, error) {
	listener, err := net.Listen("tcp", MockAddressTELNET)
	if err != nil {
		return nil, err
	}

	server := &MockServerTELNET{
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
func (s *MockServerTELNET) Close() error {
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
func (s *MockServerTELNET) Addr() string {
	return s.addr
}

// serve handles incoming requests until a stop signal is given with Close.
func (s *MockServerTELNET) serve() {
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
func (s *MockServerTELNET) handle(conn net.Conn) {
	s.mu.Lock()
	s.connections[conn] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.closeConnection(conn)
		s.wg.Done()
	}()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	defer w.Flush()

	if !s.auth(r, w) {
		return
	}

	scanner := bufio.NewScanner(r)
	for {
		scanned := scanner.Scan()
		if !scanned {
			if err := scanner.Err(); err != nil {
				if err == io.EOF {
					return
				}

				s.reportError(fmt.Errorf("handle read request error: %s", err))
				return
			}

			break
		}

		request := scanner.Text()

		switch request {
		case "":
		case MockCommandHelpTELNET:
			w.WriteString(MockCommandHelpResponseTELNET + telnet.CRLF)
		case "exit":
		default:
			w.WriteString(fmt.Sprintf("*** ERROR: unknown command '%s'", request) + telnet.CRLF)
		}

		w.Flush()
	}
}

// isRunning returns true if MockServer is running and false if is not.
func (s *MockServerTELNET) isRunning() bool {
	select {
	case <-s.quit:
		return false
	default:
		return true
	}
}

// closeConnection closes a client conn and removes it from connections map.
func (s *MockServerTELNET) closeConnection(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := conn.Close(); err != nil {
		s.reportError(fmt.Errorf("close conn error: %s", err))
	}
	delete(s.connections, conn)
}

// reportError writes error to errors channel.
func (s *MockServerTELNET) reportError(err error) bool {
	if err == nil {
		return false
	}

	select {
	case s.errors <- err:
		return true
	default:
		fmt.Printf("erros channel is locked: %s\n", err)
		// panic("erros channel is locked")
		return false
	}
}

// auth checks authorisation data and returns true if received password is valid.
func (s *MockServerTELNET) auth(r *bufio.Reader, w *bufio.Writer) bool {
	const limit = 10

	w.WriteString("Please enter password:" + telnet.CRLF)
	defer w.Flush()

	for attempt := 1; attempt <= limit; attempt++ {
		w.Flush()

		p := make([]byte, len([]byte(MockPasswordTELNET)))
		r.Read(p)
		password := string(p)

		switch password {
		case MockPasswordTELNET:
			w.WriteString(telnet.AuthSuccess + telnet.CRLF)
			return true
		default:
			if attempt == limit {
				w.WriteString(telnet.AuthTooManyFails + telnet.CRLF)
				return false
			}

			w.WriteString(telnet.AuthIncorrectPassword + telnet.CRLF)
		}
	}

	return false
}
