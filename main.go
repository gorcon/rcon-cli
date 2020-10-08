package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gorcon/rcon-cli/internal/config"
	"github.com/gorcon/rcon-cli/internal/proto/rcon"
	"github.com/gorcon/rcon-cli/internal/proto/telnet"
	"github.com/gorcon/rcon-cli/internal/session"
	"github.com/urfave/cli"
)

// Defaults.
const (
	// DefaultConfigName sets the default config file name.
	DefaultConfigName = "rcon.yaml"

	// DefaultConfigEnv is the name of the environment, which is taken
	// as default unless another value is passed.
	DefaultConfigEnv = "default"

	// DefaultLogName sets the default log file name.
	DefaultLogName = "rcon-default.log"
)

// CommandQuit is the command for exit from Interactive mode.
const CommandQuit = ":q"

// LogRecordTimeLayout is layout for convert time.Now to String.
const LogRecordTimeLayout = "2006-01-02 15:04:05"

// LogRecordFormat is format to log line record.
const LogRecordFormat = "[%s] %s: %s\n%s\n\n"

// LogFileName is the name of the file to which requests will be logged.
// If not specified, no logging will be performed.
// TODO: replace global LogFileName to better implementation.
var LogFileName string

// Version displays service version in semantic versioning (http://semver.org/).
// Can be replaced while compiling with flag `-ldflags "-X main.Version=${VERSION}"`.
var Version = "develop"

func main() {
	app := NewApp(os.Stdin, os.Stdout)

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// NewApp creates a new cli Application.
func NewApp(r io.Reader, w io.Writer) *cli.App {
	app := cli.NewApp()
	app.Usage = "CLI for executing queries on a remote server"
	app.Description = "Can be run in two modes - in the mode of a single query" +
		"\n   and in the mode of reading the input stream"
	app.Version = Version
	app.Copyright = "Copyright (c) 2020 Pavel Korotkiy (outdead)"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "a, address",
			Usage: "Set host and port to remote server. Example 127.0.0.1:16260" +
				"\n                              can be set in the config file " + DefaultConfigName + ".",
		},
		cli.StringFlag{
			Name: "p, password",
			Usage: "Set password to remote server" +
				"\n                               can be set in the config file " + DefaultConfigName + ".",
		},
		cli.StringFlag{
			Name:  "c, command",
			Usage: "Command to execute on remote server. Required flag to run in single mode",
		},
		cli.StringFlag{
			Name: "e, env",
			Usage: "Allows to select remote server address and password from the environment" +
				"\n                              in the configuration file",
		},
		cli.StringFlag{
			Name:  "l, log",
			Usage: "Path and name of the log file. if not specified, it is taken from the config.",
		},
		cli.StringFlag{
			Name: "cfg",
			Usage: "Allows to specify the path and name of the configuration file. The default" +
				"\n                value is " + DefaultConfigName + ".",
		},
		cli.StringFlag{
			Name:  "t, type",
			Usage: "Allows to specify type of connection. The default value is " + session.DefaultProtocol + ".",
		},
	}
	app.Action = func(c *cli.Context) error {
		ses, err := GetCredentials(c)
		if err != nil {
			return err
		}

		command := c.String("command")
		if command == "" {
			return Interactive(r, w, ses)
		}

		if ses.Address == "" {
			return errors.New("address is not set: to set address add -a host:port")
		}

		if ses.Password == "" {
			return errors.New("password is not set: to set password add -p password")
		}

		return Execute(w, ses, command)
	}

	return app
}

// Execute sends command to Execute to the remote server and prints the response.
func Execute(w io.Writer, ses session.Session, command string) error {
	if command == "" {
		return errors.New("command is not set")
	}

	var result string
	var err error

	switch ses.Type {
	case session.ProtocolTELNET:
		result, err = telnet.Execute(ses.Address, ses.Password, command)
	default:
		result, err = rcon.Execute(ses.Address, ses.Password, command)
	}

	if result != "" {
		fmt.Fprintln(w, result)
	}

	if err != nil {
		return err
	}

	if err := AddLog(LogFileName, ses.Address, command, result); err != nil {
		return fmt.Errorf("log error: %s", err)
	}

	return nil
}

// Interactive reads stdin, parses commands, executes them on remote server
// and prints the responses.
func Interactive(r io.Reader, w io.Writer, ses session.Session) error {
	if ses.Address == "" {
		fmt.Fprint(w, "Enter remote host and port [ip:port]: ")
		fmt.Fscanln(r, &ses.Address)
	}

	switch ses.Type {
	case session.ProtocolTELNET:
		return telnet.Interactive(r, w, ses.Address, ses.Password)
	default:
		// Default type is RCON.
		if ses.Password == "" {
			fmt.Fprint(w, "Enter password: ")
			fmt.Fscanln(r, &ses.Password)
		}

		if err := rcon.CheckCredentials(ses.Address, ses.Password); err != nil {
			return err
		}

		fmt.Fprintf(w, "Waiting commands for %s (or type %s to exit)\n> ", ses.Address, CommandQuit)

		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			command := scanner.Text()
			if command != "" {
				if command == CommandQuit {
					break
				}

				if err := Execute(w, ses, command); err != nil {
					return err
				}
			}

			fmt.Fprint(w, "> ")
		}
	}

	return nil
}

// GetCredentials parses os args or config file for details of connecting to
// a remote server. If the address and password flags were received, the
// configuration file is ignored.
func GetCredentials(c *cli.Context) (ses session.Session, err error) {
	ses.Address = c.GlobalString("a")
	ses.Password = c.GlobalString("p")
	LogFileName = c.GlobalString("l")
	ses.Type = c.GlobalString("t")

	if ses.Address != "" && ses.Password != "" {
		return ses, nil
	}

	cfg, err := GetConfig(c)
	if err != nil {
		return ses, err
	}

	e := c.GlobalString("e")
	if e == "" {
		e = DefaultConfigEnv
	}

	// Get address from environment in config if -a flag is not defined.
	if ses.Address == "" {
		ses.Address = (*cfg)[e].Address
	}

	// Get password from environment in config if -p flag is not defined.
	if ses.Password == "" {
		ses.Password = (*cfg)[e].Password
	}

	if LogFileName == "" {
		LogFileName = (*cfg)[e].Log
	}

	if ses.Type == "" {
		ses.Type = (*cfg)[e].Type
	}

	return ses, err
}

// GetConfig finds and parses config file for details of connecting to
// a remote server.
func GetConfig(c *cli.Context) (*config.Config, error) {
	path := c.GlobalString("cfg")
	if path == "" {
		home, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return nil, err
		}

		path = home + "/" + DefaultConfigName
		if _, err := os.Stat(path); err != nil {
			return nil, err
		}
	}

	// Read the config file if file exists.
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	cfg, err := config.ReadYamlConfig(path)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// AddLog saves request and response to log file.
func AddLog(logName string, address string, request string, response string) error {
	// Disable logging if log file name is empty.
	if logName == "" {
		return nil
	}

	file, err := GetLogFile(logName)
	if err != nil {
		return err
	}
	defer file.Close()

	line := fmt.Sprintf(LogRecordFormat, time.Now().Format(LogRecordTimeLayout), address, request, response)
	if _, err := file.WriteString(line); err != nil {
		return err
	}

	return nil
}

// GetLogFile opens file for append strings. Creates file if file not exist.
func GetLogFile(logName string) (*os.File, error) {
	if logName == "" {
		return nil, errors.New("empty file name")
	}

	var file *os.File

	_, err := os.Stat(logName)

	switch {
	case err == nil:
		// Open current file.
		file, err = os.OpenFile(logName, os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			return nil, err
		}
	case os.IsNotExist(err):
		// Create new file.
		file, err = os.Create(logName)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	return file, nil
}
