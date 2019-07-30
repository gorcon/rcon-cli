package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/gorcon/rcon"
	"github.com/urfave/cli"
)

// DefaultConfigName sets the default config file name.
const DefaultConfigName = "rcon.yaml"

// DefaultConfigEnv is the name of the environment, which is taken
// as default unless another value is passed.
const DefaultConfigEnv = "default"

// CommandQuit is the command for exit from interactive mode.
const CommandQuit = ":q"

// LogRecordTimeLayout is layout for convert time.Now to String
const LogRecordTimeLayout = "2006-01-02 15:04:05"

// LogRecordFormat is format to log line record.
const LogRecordFormat = "[%s] %s: %s\n%s\n\n"

// LogFileName is the name of the file to which requests will be logged.
// If not specified, no logging will be performed.
var LogFileName string

// Config allows to take a remote server address and password from
// the configuration file. This enables not to specify these flags when
// running the CLI.
//
// Example:
// ```yaml
// default:
//   address: "127.0.0.1:16260"
//   password: "password"
// ```
type Config map[string]struct {
	Address  string `json:"address" yaml:"address"`
	Password string `json:"password" yaml:"password"`
	Log      string `json:"log" yaml:"log"`
}

func main() {
	var description = "Can be run in two modes - in the mode of a single query"
	description += "\n   and in the mode of reading the input stream"

	app := cli.NewApp()
	app.Usage = "CLI for executing queries on a remote server"
	app.Description = description
	app.Version = "0.3.1"
	app.Copyright = "Copyright (c) 2019 Pavel Korotkiy (outdead)"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "a, address",
			Usage: "set host and port to remote rcon server. Example 127.0.0.1:16260" +
				"\n                              can be set in the config file rcon.yaml",
		},
		cli.StringFlag{
			Name: "p, password",
			Usage: "set password to remote rcon server" +
				"\n                               can be set in the config file rcon.yaml",
		},
		cli.StringFlag{
			Name:  "c, command",
			Usage: "command to execute on remote server. Required flag to run in single mode",
		},
		cli.StringFlag{
			Name: "e, env",
			Usage: "allows to select remote server address and password from the environment" +
				"\n                              in the configuration file",
		},
		cli.StringFlag{
			Name:  "l, log",
			Usage: "path and name of the log file. if not specified, it is taken from the config.",
		},
		cli.StringFlag{
			Name: "cfg",
			Usage: "allows to specify the path and name of the configuration file. The default" +
				"\n                value is rcon.yaml.",
		},
	}
	app.Action = func(c *cli.Context) error {
		address, password, err := getCredentials(c)
		if err != nil {
			return err
		}

		command := c.String("command")
		if command == "" {
			return interactive(address, password)
		}

		if address == "" || password == "" {
			if address == "" {
				return errors.New("address is not set: to set address add -a host:port")
			}

			if password == "" {
				return errors.New("password is not set: to set password add -p password")
			}
		}

		return execute(address, password, command)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// execute sends command to execute to the remote server and prints the response.
func execute(address string, password string, command string) error {
	console, err := rcon.Dial(address, password)
	if err != nil {
		return err
	}
	defer console.Close()

	result, err := console.Execute(command)
	if result != "" {
		fmt.Println(result)
	}

	if err := addLog(LogFileName, address, command, result); err != nil {
		err = fmt.Errorf("log error: %s", err)
		fmt.Println(err)
	}

	return err
}

// interactive reads stdin, parses commands, executes them on remote server
// and prints the responses.
func interactive(address string, password string) error {
	if address == "" {
		fmt.Print("enter host and port from remote server: ")
		fmt.Scanln(&address)
	}

	if password == "" {
		fmt.Print("enter the password: ")
		fmt.Scanln(&password)
	}

	if err := checkCredentials(address, password); err != nil {
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("waiting commands for %s\n> ", address)
	for scanner.Scan() {
		command := scanner.Text()
		if command != "" {
			if command == CommandQuit {
				return nil
			}

			if err := execute(address, password, command); err != nil {
				return err
			}
		}

		fmt.Print("> ")
	}

	return nil
}

// readYamlConfig reads config data from yaml file.
func readYamlConfig(path string) (Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// getCredentials parses os args or config file for details of connecting to
// a remote server. If the address and password flags were received, the
// configuration file is ignored.
func getCredentials(c *cli.Context) (address string, password string, err error) {
	address = c.GlobalString("a")
	password = c.GlobalString("p")
	LogFileName = c.GlobalString("l")

	if address != "" && password != "" {
		return address, password, nil
	}

	path := c.GlobalString("cfg")
	if path == "" {
		var home string
		home, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return address, password, err
		}
		path = home + "/" + DefaultConfigName

		if _, err2 := os.Stat(path); err2 != nil {
			// TODO: Think about creation of the config file.
			return address, password, err
		}
	}

	// Read the config file if file exists.
	_, err = os.Stat(path)
	if err == nil {
		cfg, err := readYamlConfig(path)
		if err != nil {
			return address, password, err
		}

		e := c.GlobalString("e")
		if e == "" {
			e = DefaultConfigEnv
		}

		// Get address from environment in config if -a flag is not defined.
		if address == "" {
			address = cfg[e].Address
		}

		// Get password from environment in config if -p flag is not defined.
		if password == "" {
			password = cfg[e].Password
		}

		if LogFileName == "" {
			LogFileName = cfg[e].Log
		}
	}

	return
}

// checkCredentials sends auth request for remote server. Returns en error if
// address or password is incorrect.
func checkCredentials(address string, password string) error {
	console, err := rcon.Dial(address, password)
	if err != nil {
		return err
	}

	return console.Close()
}

// addLog saves request and response to log file.
func addLog(logName string, address string, request string, response string) error {
	if logName == "" {
		return nil
	}

	var file *os.File
	if _, err := os.Stat(logName); os.IsNotExist(err) {
		file, err = os.Create(logName)
		if err != nil {
			return err
		}
	} else {
		file, err = os.OpenFile(logName, os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	now := time.Now()
	line := fmt.Sprintf(LogRecordFormat, now.Format(LogRecordTimeLayout), address, request, response)
	if _, err := file.WriteString(line); err != nil {
		return err
	}

	return nil
}
