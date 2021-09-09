package executor

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorcon/rcon"
	"github.com/gorcon/rcon-cli/internal/config"
	"github.com/gorcon/rcon-cli/internal/logger"
	"github.com/gorcon/telnet"
	"github.com/gorcon/websocket"
	"github.com/urfave/cli/v2"
)

// CommandQuit is the command for exit from Interactive mode.
const CommandQuit = ":q"

// AttemptsLimit is the limit value for the number of attempts to obtain user
// data in terminal mode.
const AttemptsLimit = 3

// CommandsResponseSeparator is symbols that is written between responses of
// several commands if more than one command was called.
// TODO: Add to config.
const CommandsResponseSeparator = "--------"

// Errors.
var (
	// ErrEmptyAddress is returned when executed command without setting address
	// in single mode.
	ErrEmptyAddress = errors.New("address is not set: to set address add -a host:port")

	// ErrEmptyPassword is returned when executed command without setting password
	// in single mode.
	ErrEmptyPassword = errors.New("password is not set: to set password add -p password")

	// ErrToManyFails is returned in terminal mode when exceeding the limit of
	// user data retrieval attempts.
	ErrToManyFails = errors.New("to many fails")

	// ErrCommandEmpty is returned when executed command length equal 0.
	ErrCommandEmpty = errors.New("command is not set")
)

// Executor is a cli commands execute wrapper.
type Executor struct {
	version string
	r       io.Reader
	w       io.Writer
	app     *cli.App

	client interface {
		Execute(command string) (string, error)
		Close() error
	}
}

// NewExecutor creates a new Executor.
func NewExecutor(r io.Reader, w io.Writer, version string) *Executor {
	executor := Executor{
		version: version,
		r:       r,
		w:       w,
	}

	return &executor
}

// Run is the entry point to the cli app.
func (executor *Executor) Run(arguments []string) error {
	executor.init()
	defer executor.Close()

	if err := executor.app.Run(arguments); err != nil && !errors.Is(err, flag.ErrHelp) {
		return fmt.Errorf("cli: %w", err)
	}

	return nil
}

// NewSession parses os args and config file for connection details to
// a remote server. If the address and password flags were received the
// configuration file is ignored.
func (executor *Executor) NewSession(c *cli.Context) (*config.Session, error) {
	ses := config.Session{
		Address:    c.String("address"),
		Password:   c.String("password"),
		Type:       c.String("type"),
		Log:        c.String("log"),
		SkipErrors: c.Bool("skip"),
	}

	if ses.Address != "" && ses.Password != "" {
		return &ses, nil
	}

	cfg, err := config.NewConfig(c.String("config"))
	if err != nil {
		return &ses, fmt.Errorf("config: %w", err)
	}

	e := c.String("env")
	if e == "" {
		e = config.DefaultConfigEnv
	}

	// Get variables from config environment if flags are not defined.
	if ses.Address == "" {
		ses.Address = (*cfg)[e].Address
	}

	if ses.Password == "" {
		ses.Password = (*cfg)[e].Password
	}

	if ses.Log == "" {
		ses.Log = (*cfg)[e].Log
	}

	if ses.Type == "" {
		ses.Type = (*cfg)[e].Type
	}

	return &ses, nil
}

// init creates a new cli Application.
func (executor *Executor) init() {
	app := cli.NewApp()
	app.Usage = "CLI for executing queries on a remote server"
	app.Description = "Can be run in two modes - in the mode of a single query and in terminal mode of reading the " +
		"input stream. \n\n" + "To run single mode type commands after options flags. Example: \n" +
		filepath.Base(os.Args[0]) + " -a 127.0.0.1:16260 -p password command1 command2 \n\n" +
		"To run terminal mode just do not specify commands to execute. Example: \n" +
		filepath.Base(os.Args[0]) + " -a 127.0.0.1:16260 -p password"
	app.Version = executor.version
	app.Copyright = "Copyright (c) 2020 Pavel Korotkiy (outdead)"
	app.HideHelpCommand = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"a"},
			Usage:   "Set host and port to remote server. Example 127.0.0.1:16260",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"p"},
			Usage:   "Set password to remote server",
		},
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Usage:   "Specify type of connection (default: " + config.DefaultProtocol + ")",
		},
		&cli.StringFlag{
			Name:    "log",
			Aliases: []string{"l"},
			Usage:   "Path to the log file. If not specified it is taken from the config",
		},
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to the configuration file (default: " + config.DefaultConfigName + ")",
		},
		&cli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "Config environment with server credentials (default: " + config.DefaultConfigEnv + ")",
		},
		&cli.BoolFlag{
			Name:    "skip",
			Aliases: []string{"s"},
			Usage:   "Skip errors and run next command",
		},
	}
	app.Action = func(c *cli.Context) error {
		ses, err := executor.NewSession(c)
		if err != nil {
			return err
		}

		commands := c.Args().Slice()
		if len(commands) == 0 {
			return executor.Interactive(executor.r, executor.w, ses)
		}

		if ses.Address == "" {
			return ErrEmptyAddress
		}

		if ses.Password == "" {
			return ErrEmptyPassword
		}

		return executor.Execute(executor.w, ses, commands...)
	}

	executor.app = app
}

// Execute sends command to Execute to the remote server and prints the response.
func (executor *Executor) Execute(w io.Writer, ses *config.Session, commands ...string) error {
	if len(commands) == 0 {
		return ErrCommandEmpty
	}

	err := executor.Dial(ses)

	if ses.Type == config.ProtocolWebRCON {
		defer func() {
			executor.client.Close()
			executor.client = nil
		}()
	}

	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}

	for i, command := range commands {
		if command == "" {
			return ErrCommandEmpty
		}

		var result string
		var err error

		result, err = executor.client.Execute(command)
		if result != "" {
			result = strings.TrimSpace(result)
			fmt.Fprintln(w, result)
		}

		if err != nil {
			if ses.SkipErrors {
				fmt.Fprintln(w, fmt.Errorf("execute: %w", err))
			} else {
				return fmt.Errorf("execute: %w", err)
			}
		}

		if err := logger.Write(ses.Log, ses.Address, command, result); err != nil {
			fmt.Fprintln(w, fmt.Errorf("log: %w", err))
		}

		if i+1 != len(commands) {
			fmt.Fprintln(w, CommandsResponseSeparator)
		}
	}

	return nil
}

// Interactive reads stdin, parses commands, executes them on remote server
// and prints the responses.
func (executor *Executor) Interactive(r io.Reader, w io.Writer, ses *config.Session) error {
	if ses.Address == "" {
		fmt.Fprint(w, "Enter remote host and port [ip:port]: ")
		fmt.Fscanln(r, &ses.Address)
	}

	if ses.Password == "" {
		fmt.Fprint(w, "Enter password: ")
		fmt.Fscanln(r, &ses.Password)
	}

	if ses.Type == "" {
		fmt.Fprint(w, "Enter protocol type (empty for rcon): ")
		fmt.Fscanln(r, &ses.Type)
	}

	switch ses.Type {
	case config.ProtocolTELNET:
		return telnet.DialInteractive(r, w, ses.Address, ses.Password)
	case "", config.ProtocolRCON, config.ProtocolWebRCON:
		if err := executor.Dial(ses); err != nil {
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

				if err := executor.Execute(w, ses, command); err != nil {
					return err
				}
			}

			fmt.Fprint(w, "> ")
		}
	default:
		fmt.Fprintf(w, "Unsupported protocol type (%q). Allowed %q, %q and %q protocols\n",
			ses.Type, config.ProtocolRCON, config.ProtocolWebRCON, config.ProtocolTELNET)
	}

	return nil
}

func (executor *Executor) Close() error {
	if executor.client != nil {
		return executor.client.Close()
	}

	return nil
}

// Dial sends auth request for remote server. Returns en error if
// address or password is incorrect.
func (executor *Executor) Dial(ses *config.Session) error {
	var err error

	if executor.client == nil {
		switch ses.Type {
		case config.ProtocolTELNET:
			executor.client, err = telnet.Dial(ses.Address, ses.Password)
		case config.ProtocolWebRCON:
			executor.client, err = websocket.Dial(ses.Address, ses.Password)
		default:
			executor.client, err = rcon.Dial(ses.Address, ses.Password)
		}
	}

	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	return nil
}
