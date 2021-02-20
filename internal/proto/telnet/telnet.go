package telnet

import (
	"fmt"
	"io"

	"github.com/gorcon/telnet"
)

// Execute sends command to Execute to the remote server and returns
// the response.
func Execute(address string, password string, command string) (string, error) {
	if command == "" {
		return "", telnet.ErrCommandEmpty
	}

	console, err := telnet.Dial(address, password)
	if err != nil {
		return "", fmt.Errorf("telnet: %w", err)
	}
	defer console.Close()

	return console.Execute(command)
}

// Interactive parses commands from input reader, executes them on remote
// server and writes responses to output writer. Password can be empty string.
// In this case password will be prompted in an interactive window.
func Interactive(r io.Reader, w io.Writer, address string, password string) error {
	return telnet.DialInteractive(r, w, address, password)
}

// CheckCredentials sends auth request for remote server. Returns en error if
// address or password is incorrect.
func CheckCredentials(address string, password string) error {
	console, err := telnet.Dial(address, password)
	if err != nil {
		return fmt.Errorf("telnet: %w", err)
	}

	return console.Close()
}
