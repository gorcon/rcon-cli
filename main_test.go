package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/gorcon/rcon-cli/internal/session"

	"github.com/gorcon/rcon-cli/internal/config"
	"github.com/stretchr/testify/assert"
)

func newConfig() *config.Config {
	return &config.Config{
		"default": session.Session{Address: "", Password: "", Log: "rcon-default.log"},
	}
}

func TestReadYamlConfig(t *testing.T) {
	func() {
		cfg, err := config.ReadYamlConfig("rcon.yaml")
		assert.NoError(t, err)
		assert.Equal(t, newConfig(), &cfg)
	}()

	func() {
		cfg, err := config.ReadYamlConfig("nonexist.yaml")
		assert.NotNil(t, err)
		var expected config.Config
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
		assert.NoError(t, err)
	}()

	// Test skip log. No logs is available.
	func() {
		err := AddLog("", address, command, result)
		assert.NoError(t, err)
	}()

	// Test create file log.
	func() {
		err := AddLog(logName, address, command, result)
		assert.NoError(t, err)
	}()

	// Test append to log file.
	func() {
		err := AddLog(logName, address, command, result)
		assert.NoError(t, err)
	}()
}

func TestGetLogFile(t *testing.T) {
	logDir := "temp"
	logName := "tmpfile.log"
	logPath := logDir + "/" + logName

	// Test empty log file name.
	func() {
		file, err := GetLogFile("")
		assert.Nil(t, file)
		assert.EqualError(t, err, "empty file name")
	}()

	// Test stat permission denied.
	func() {
		if err := os.Mkdir(logDir, 0400); err != nil {
			assert.NoError(t, err)
			return
		}
		defer func() {
			err := os.RemoveAll(logDir)
			assert.NoError(t, err)
		}()

		file, err := GetLogFile(logPath)
		assert.Nil(t, file)
		assert.EqualError(t, err, fmt.Sprintf("stat %s: permission denied", logPath))
	}()

	// Test create permission denied.
	func() {
		if err := os.Mkdir(logDir, 0500); err != nil {
			assert.NoError(t, err)
			return
		}
		defer func() {
			err := os.RemoveAll(logDir)
			assert.NoError(t, err)
		}()

		file, err := GetLogFile(logPath)
		assert.Nil(t, file)
		assert.EqualError(t, err, fmt.Sprintf("open %s: permission denied", logPath))
	}()

	// Positive test create new log file + test open permission denied.
	func() {
		if err := os.Mkdir(logDir, 0700); err != nil {
			assert.NoError(t, err)
			return
		}
		defer func() {
			err := os.RemoveAll(logDir)
			assert.NoError(t, err)
		}()

		// Positive test create new log file.
		file, err := GetLogFile(logPath)
		assert.NotNil(t, file)
		assert.NoError(t, err)

		if err := os.Chmod(logPath, 0000); err != nil {
			assert.NoError(t, err)
			return
		}

		// Test open permission denied.
		file, err = GetLogFile(logPath)
		assert.Nil(t, file)
		assert.EqualError(t, err, fmt.Sprintf("open %s: permission denied", logPath))
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
		err := Execute(w, session.Session{Address: "", Password: MockPassword}, MockCommandHelp)
		assert.Error(t, err)
	}()

	// Test empty password.
	func() {
		err := Execute(w, session.Session{Address: server.Addr(), Password: ""}, MockCommandHelp)
		assert.Error(t, err)
	}()

	// Test wrong password.
	func() {
		err := Execute(w, session.Session{Address: server.Addr(), Password: "wrong"}, MockCommandHelp)
		assert.Error(t, err)
	}()

	// Test empty command.
	func() {
		err := Execute(w, session.Session{Address: server.Addr(), Password: MockPassword}, "")
		assert.Error(t, err)
	}()

	// Test long command.
	func() {
		bigCommand := make([]byte, 1001)
		err := Execute(w, session.Session{Address: server.Addr(), Password: MockPassword}, string(bigCommand))
		assert.Error(t, err)
	}()

	// Positive test Execute func.
	func() {
		err := Execute(w, session.Session{Address: server.Addr(), Password: MockPassword}, MockCommandHelp)
		assert.NoError(t, err)
	}()

	// Positive test Execute func with log.
	func() {
		LogFileName = "tmpfile.log"
		defer func() {
			err := os.Remove(LogFileName)
			assert.NoError(t, err)
			LogFileName = ""
		}()

		err := Execute(w, session.Session{Address: server.Addr(), Password: MockPassword}, MockCommandHelp)
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

		err = Interactive(&r, w, session.Session{Address: server.Addr(), Password: "fake"})
		assert.Error(t, err)
	}()

	// Test get Interactive address.
	func() {
		var r bytes.Buffer
		r.WriteString(server.Addr() + "\n")
		r.WriteString(CommandQuit + "\n")

		err = Interactive(&r, w, session.Session{Address: "", Password: MockPassword})
		assert.NoError(t, err)
	}()

	// Test get Interactive password.
	func() {
		var r bytes.Buffer
		r.WriteString(MockPassword + "\n")
		r.WriteString(CommandQuit + "\n")

		err = Interactive(&r, w, session.Session{Address: server.Addr(), Password: ""})
		assert.NoError(t, err)
	}()

	// Test get Interactive commands.
	func() {
		r := &bytes.Buffer{}
		r.WriteString(MockCommandHelp + "\n")
		r.WriteString("unknown command" + "\n")
		r.WriteString(CommandQuit + "\n")

		err = Interactive(r, w, session.Session{Address: server.Addr(), Password: MockPassword})
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
		assert.NoError(t, err)
		defer func() {
			err := os.Remove(configFileName)
			assert.NoError(t, err)
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

	// Test default config file not exist. Log is not used.
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

	// Test default config file is incorrect. Log is not used.
	func() {
		var configFileName = "rcon-temp.yaml"
		err := CreateInvalidConfigFile(configFileName, server.Addr(), MockPassword)
		assert.NoError(t, err)
		defer func() {
			err := os.Remove(configFileName)
			assert.NoError(t, err)
		}()

		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		args = append(args, "-cfg="+configFileName)
		args = append(args, "-c="+MockCommandHelp)

		err = app.Run(args)
		assert.EqualError(t, err, "yaml: line 1: did not find expected key")
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

	// Test empty password. Log is not used.
	func() {
		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		// Hack to use os.Args[0] in go run
		args[0] = ""
		args = append(args, "-a="+server.Addr())
		args = append(args, "-c="+MockCommandHelp)

		err = app.Run(args)
		assert.EqualError(t, err, "password is not set: to set password add -p password")
	}()

	// Positive test Interactive. Log is not used.
	func() {
		r := &bytes.Buffer{}
		w := &bytes.Buffer{}

		app := NewApp(r, w)
		args := os.Args[0:1]
		args = append(args, "-a="+server.Addr())
		args = append(args, "-p="+MockPassword)

		r.WriteString(MockCommandHelp + "\n")
		r.WriteString(CommandQuit + "\n")

		err = app.Run(args)
		assert.NoError(t, err)
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

// CreateIncorrectConfigFile creates incorrect yaml config file.
func CreateInvalidConfigFile(name string, address string, password string) error {
	var stringBody = fmt.Sprintf(
		"address: \"%s\"\n  password: \"%s\"\n  log: \"%s\"",
		address, password, DefaultLogName,
	)
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	_, err = file.WriteString(stringBody)

	return err
}
