package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
	"github.com/gorcon/rcon"
	"github.com/urfave/cli"
)

// DefaultConfigName sets the default config file name.
const DefaultConfigName = "rcon.yaml"

// DefaultConfigEnv is the name of the environment, which is taken
// as default unless another value is passed.
const DefaultConfigEnv = "default"

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
}

func main() {
	var description = "Can be run in two modes - in the mode of a single query"
	description += "\n   and in the mode of reading the input stream"

	app := cli.NewApp()
	app.Usage = "CLI for executing queries on a remote server"
	app.Description = description
	app.Version = "0.2.0"
	app.Author = "Pavel Korotkiy (outdead)"
	app.Copyright = "Copyright (c) 2019 Pavel Korotkiy"
	app.Commands = []cli.Command{
		{
			Name:  "cli",
			Usage: "Run CLI for commands in the form of successive lines of text from the input stream",
			Action: func(c *cli.Context) error {
				address, password := getCredentials(c)
				if address == "" || password == "" {
					cli.ShowAppHelpAndExit(c, 1)
					return nil
				}

				scanner := bufio.NewScanner(os.Stdin)
				fmt.Printf("waiting commands for %s\n", address)
				fmt.Print("> ")
				for scanner.Scan() {
					command := scanner.Text()
					if command != "" {
						if command == ":q" {
							return nil
						}

						execute(address, password, command)
					}

					fmt.Print("> ")
				}

				return nil
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "address, a",
			Usage: "set host and port to remote rcon server. Example 127.0.0.1:16260" +
				"\n                              can be set in the config file rcon.yaml",
		},
		cli.StringFlag{
			Name: "password, p",
			Usage: "set password to remote rcon server" +
				"\n                               can be set in the config file rcon.yaml",
		},
		cli.StringFlag{
			Name:  "command, c",
			Usage: "command to execute on remote server. Required flag to run in single mode",
		},
		cli.StringFlag{
			Name: "env, e",
			Usage: "allows to select remote server address and password from the environment " +
				"\n                              in the configuration file",
		},
	}
	app.Action = func(c *cli.Context) error {
		address, password := getCredentials(c)
		command := c.String("command")
		if address == "" || password == "" || command == "" {
			cli.ShowAppHelpAndExit(c, 1)
			return nil
		}

		execute(address, password, command)
		return nil
	}

	app.Run(os.Args)
}

// execute sends command to execute to the remote server.
func execute(address string, password string, command string) {
	console, err := rcon.Dial(address, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer console.Close()

	result, err := console.Execute(command)
	if err != nil {
		fmt.Println(err)
	}

	if result != "" {
		fmt.Println(result)
	}
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
func getCredentials(c *cli.Context) (address string, password string) {
	address = c.GlobalString("a")
	password = c.GlobalString("p")

	if address != "" && password != "" {
		return
	}

	home, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return
	}
	path := home + "/" + DefaultConfigName

	if _, err := os.Stat(path); err == nil {
		cfg, err := readYamlConfig(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		e := c.GlobalString("e")
		if e == "" {
			e = DefaultConfigEnv
		}

		if address == "" {
			address = cfg[e].Address
		}

		if password == "" {
			password = cfg[e].Password
		}
	}

	return
}
