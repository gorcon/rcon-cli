package rcon

import (
	"errors"

	"github.com/gorcon/rcon"
)

// Execute sends command to Execute to the remote server and returns
// the response.
func Execute(address string, password string, command string) (string, error) {
	if command == "" {
		return "", errors.New("command is not set")
	}

	console, err := rcon.Dial(address, password)
	if err != nil {
		return "", err
	}
	defer console.Close()

	return console.Execute(command)
}

// CheckCredentials sends auth request for remote server. Returns en error if
// address or password is incorrect.
func CheckCredentials(address string, password string) error {
	console, err := rcon.Dial(address, password)
	if err != nil {
		return err
	}

	return console.Close()
}
