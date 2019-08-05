package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newConfig() *Config {
	return &Config{
		"default": struct {
			Address  string `json:"address" yaml:"address"`
			Password string `json:"password" yaml:"password"`
			Log      string `json:"log" yaml:"log"`
		}{Address: "", Password: "", Log: "rcon-default.log"},
	}
}

func TestReadYamlConfig(t *testing.T) {
	func() {
		cfg, err := ReadYamlConfig("rcon.yaml")
		assert.Nil(t, err)
		assert.Equal(t, newConfig(), &cfg)
	}()

	func() {
		cfg, err := ReadYamlConfig("nonexist.yaml")
		assert.NotNil(t, err)
		var expected Config
		assert.Equal(t, expected, cfg)
	}()
}

func TestAddLog(t *testing.T) {
	logName := "tmpfile.log"

	address := "127.0.0.1:16200"
	command := "players"
	result := `Players connected (2):
-admin
-testuser`

	defer func() {
		err := os.Remove(logName)
		assert.Nil(t, err)
	}()

	// Test skip log. No logs is available.
	func() {
		err := AddLog("", address, command, result)
		assert.Nil(t, err)
	}()

	// Test create file log.
	func() {
		err := AddLog(logName, address, command, result)
		assert.Nil(t, err)
	}()

	// Test append to log file.
	func() {
		err := AddLog(logName, address, command, result)
		assert.Nil(t, err)
	}()
}

func TestCheckCredentials(t *testing.T) {
	server, err := NewMockServer()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		assert.NoError(t, server.Close())
		close(server.errors)
		for err := range server.errors {
			assert.NoError(t, err)
		}
	}()

	// Test invalid credentials.
	func() {
		err = CheckCredentials(server.Addr(), "")
		assert.Error(t, err)
	}()

	// Positive test CheckCredentials func.
	func() {
		err = CheckCredentials(server.Addr(), MockPassword)
		assert.NoError(t, CheckCredentials(server.Addr(), MockPassword))
	}()
}

func TestExecute(t *testing.T) {
	server, err := NewMockServer()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		assert.NoError(t, server.Close())
		close(server.errors)
		for err := range server.errors {
			assert.NoError(t, err)
		}
	}()

	w := &bytes.Buffer{}

	// Test empty address.
	func() {
		err := Execute(w, "", MockPassword, MockCommandHelp)
		assert.Error(t, err)
	}()

	// Test empty password.
	func() {
		err := Execute(w, server.Addr(), "", MockCommandHelp)
		assert.Error(t, err)
	}()

	// Test wrong password.
	func() {
		err := Execute(w, server.Addr(), "wrong", MockCommandHelp)
		assert.Error(t, err)
	}()

	// Test empty command.
	func() {
		err := Execute(w, server.Addr(), MockPassword, "")
		assert.Error(t, err)
	}()

	// Test long command.
	func() {
		bigCommand := make([]byte, 1001)
		err := Execute(w, server.Addr(), MockPassword, string(bigCommand))
		assert.Error(t, err)
	}()

	// Positive test Execute func.
	func() {
		err := Execute(w, server.Addr(), MockPassword, MockCommandHelp)
		assert.NoError(t, err)
	}()

	// Positive test Execute func with log.
	func() {
		LogFileName = "tmpfile.log"
		defer func() {
			err := os.Remove(LogFileName)
			assert.Nil(t, err)
			LogFileName = ""
		}()

		err := Execute(w, server.Addr(), MockPassword, MockCommandHelp)
		assert.NoError(t, err)
	}()
}

func TestInteractive(t *testing.T) {
	server, err := NewMockServer()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		assert.NoError(t, server.Close())
		close(server.errors)
		for err := range server.errors {
			assert.NoError(t, err)
		}
	}()

	w := &bytes.Buffer{}

	// Test wrong password.
	func() {
		var r bytes.Buffer
		r.WriteString(CommandQuit + "\n")

		err = Interactive(&r, w, server.Addr(), "fake")
		assert.Error(t, err)
	}()

	// Test get Interactive address.
	func() {
		var r bytes.Buffer
		r.WriteString(server.Addr() + "\n")
		r.WriteString(CommandQuit + "\n")

		err = Interactive(&r, w, "", MockPassword)
		assert.NoError(t, err)
	}()

	// Test get Interactive password.
	func() {
		var r bytes.Buffer
		r.WriteString(MockPassword + "\n")
		r.WriteString(CommandQuit + "\n")

		err = Interactive(&r, w, server.Addr(), "")
		assert.NoError(t, err)
	}()

	// Test get Interactive commands.
	func() {
		r := &bytes.Buffer{}
		r.WriteString(MockCommandHelp + "\n")
		r.WriteString("unknown command" + "\n")
		r.WriteString(CommandQuit + "\n")

		err = Interactive(r, w, server.Addr(), MockPassword)
		assert.NoError(t, err)
	}()
}

func TestNewApp(t *testing.T) {
	server, err := NewMockServer()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		assert.NoError(t, server.Close())
		close(server.errors)
		for err := range server.errors {
			assert.NoError(t, err)
		}
	}()

	// Test getting address and password from args. Config ang log are not used.
	func() {
		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		args = append(args, "-a="+server.Addr())
		args = append(args, "-p="+MockPassword)
		args = append(args, "-c="+MockCommandHelp)

		err = app.Run(args)
		assert.NoError(t, err)
	}()

	// Test getting address and password from config. Log is not used.
	func() {
		var configFileName = "rcon-temp.yaml"
		err := CreateConfigFile(configFileName, server.Addr(), MockPassword)
		defer func() {
			err := os.Remove(configFileName)
			assert.Nil(t, err)
		}()

		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		args = append(args, "-cfg="+configFileName)
		args = append(args, "-c="+MockCommandHelp)

		err = app.Run(args)
		assert.NoError(t, err)
	}()

	// Test create default config file. Log is not used.
	func() {
		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		args = append(args, "-c="+MockCommandHelp)

		err = app.Run(args)
		assert.Error(t, err)
		if !os.IsNotExist(err) {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	// Test empty address and password. Log is not used.
	func() {
		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		// Hack to use os.Args[0] in go run
		args[0] = ""
		args = append(args, "-c="+MockCommandHelp)

		err = app.Run(args)
		assert.EqualError(t, err, "address is not set: to set address add -a host:port")
	}()
}

// CreateConfigFile creates config file with default section.
func CreateConfigFile(name string, address string, password string) error {
	var stringBody = fmt.Sprintf(
		"%s:\n  address: \"%s\"\n  password: \"%s\"\n  log: \"%s\"",
		DefaultConfigEnv, address, password, DefaultLogName,
	)
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	_, err = file.WriteString(stringBody)

	return err
}
