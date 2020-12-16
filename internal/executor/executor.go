package executor

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/gorcon/rcon-cli/internal/config"
	"github.com/gorcon/rcon-cli/internal/logger"
	"github.com/gorcon/rcon-cli/internal/proto/rcon"
	"github.com/gorcon/rcon-cli/internal/proto/telnet"
	"github.com/gorcon/rcon-cli/internal/proto/websocket"
	"github.com/gorcon/rcon-cli/internal/session"
	"github.com/urfave/cli"
)

// CommandQuit is the command for exit from Interactive mode.
const CommandQuit = ":q"

// Executor is a cli commands execute wrapper.
type Executor struct {
	version string
	r       io.Reader
	w       io.Writer
	app     *cli.App
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

	if err := executor.app.Run(arguments); err != nil && !errors.Is(err, flag.ErrHelp) {
		return err
	}

	return nil
}

// NewSession parses os args and config file for connection details to
// a remote server. If the address and password flags were received the
// configuration file is ignored.
func (executor *Executor) NewSession(c *cli.Context) (*session.Session, error) {
	ses := session.Session{
		Address:  c.GlobalString("a"),
		Password: c.GlobalString("p"),
		Log:      c.GlobalString("l"),
		Type:     c.GlobalString("t"),
	}

	if ses.Address != "" && ses.Password != "" {
		return &ses, nil
	}

	cfg, err := config.NewConfig(c.GlobalString("cfg"))
	if err != nil {
		return &ses, err
	}

	e := c.GlobalString("e")
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

	return &ses, err
}

// init creates a new cli Application.
func (executor *Executor) init() {
	app := cli.NewApp()
	app.Usage = "CLI for executing queries on a remote server"
	app.Description = "Can be run in two modes - in the mode of a single query" +
		"\n   and in terminal mode of reading the input stream. To run terminal mode" +
		"\n   just do not specify command to execute."
	app.Version = executor.version
	app.Copyright = "Copyright (c) 2020 Pavel Korotkiy (outdead)"
	app.HideHelp = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "a, address",
			Usage: "Set host and port to remote server. Example 127.0.0.1:16260",
		},
		cli.StringFlag{
			Name:  "p, password",
			Usage: "Set password to remote server",
		},
		cli.StringFlag{
			Name:  "t, type",
			Usage: "Allows to specify type of connection. Default value is " + session.DefaultProtocol,
		},
		cli.StringFlag{
			Name:  "l, log",
			Usage: "Path and name of the log file. If not specified, it is taken from the config",
		},
		cli.StringFlag{
			Name:  "c, command",
			Usage: "Command to execute on remote server. Required flag to run in single mode",
		},
		cli.StringFlag{
			Name:  "e, env",
			Usage: "Allows to select server credentials from selected environment in the configuration file",
		},
		cli.StringFlag{
			Name:  "cfg",
			Usage: "Allows to specify the path and name of the configuration file. Default value is " + config.DefaultConfigName,
		},
	}
	app.Action = func(c *cli.Context) error {
		ses, err := executor.NewSession(c)
		if err != nil {
			return err
		}

		command := c.String("command")
		if command == "" {
			return Interactive(executor.r, executor.w, ses)
		}

		if ses.Address == "" {
			return errors.New("address is not set: to set address add -a host:port")
		}

		if ses.Password == "" {
			return errors.New("password is not set: to set password add -p password")
		}

		return Execute(executor.w, ses, command)
	}

	executor.app = app
}

// Execute sends command to Execute to the remote server and prints the response.
func Execute(w io.Writer, ses *session.Session, command string) error {
	if command == "" {
		return errors.New("command is not set")
	}

	var result string
	var err error

	switch ses.Type {
	case session.ProtocolTELNET:
		result, err = telnet.Execute(ses.Address, ses.Password, command)
	case session.ProtocolWebRCON:
		result, err = websocket.Execute(ses.Address, ses.Password, command)
	default:
		result, err = rcon.Execute(ses.Address, ses.Password, command)
	}

	if result != "" {
		result = strings.TrimSpace(result)
		fmt.Fprintln(w, result)
	}

	if err != nil {
		return err
	}

	if err := logger.AddLog(ses.Log, ses.Address, command, result); err != nil {
		return fmt.Errorf("log error: %s", err)
	}

	return nil
}

// Interactive reads stdin, parses commands, executes them on remote server
// and prints the responses.
func Interactive(r io.Reader, w io.Writer, ses *session.Session) error {
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

		if err := CheckCredentials(ses); err != nil {
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

// CheckCredentials sends auth request for remote server. Returns en error if
// address or password is incorrect.
func CheckCredentials(ses *session.Session) error {
	if ses.Type == session.ProtocolWebRCON {
		return websocket.CheckCredentials(ses.Address, ses.Password)
	}

	return rcon.CheckCredentials(ses.Address, ses.Password)
}
