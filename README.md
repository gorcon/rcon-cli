# rcon-cli
[![Go Report Card](https://goreportcard.com/badge/github.com/gorcon/rcon-cli)](https://goreportcard.com/report/github.com/gorcon/rcon-cli)
![GitHub All Releases](https://img.shields.io/github/downloads/gorcon/rcon-cli/total)

CLI for executing queries on a remote server

## Supported Games

* [Project Zomboid](https://store.steampowered.com/app/108600) 
* [Conan Exiles](https://store.steampowered.com/app/440900)
* [Rust](https://store.steampowered.com/app/252490) (add +rcon.web 0 to the args when starting the server)
* [ARK: Survival Evolved](https://store.steampowered.com/app/346110)

Open pull request if you have successfully used a package with another game with rcon support and add it to the list.

## Installation

Download the binary for your platform from the [releases](https://github.com/gorcon/rcon-cli/releases)

## Usage

```text
USAGE:
   rcon-cli [global options] command [command options] [arguments...]

COMMANDS:
     cli      Run CLI for commands in the form of successive lines of text from the input stream
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --address value, -a value  set host and port to remote rcon server. Example 127.0.0.1:16260
                              can be set in the config file rcon.yaml
   --password value, -p value  set password to remote rcon server
                               can be set in the config file rcon.yaml
   --command value, -c value  command to execute on remote server. Required flag to run in single mode
   --env value, -e value      allows to select remote server address and password from the environment 
                              in the configuration file
   --help, -h     show help
   --version, -v  print the version
```

Rcon CLI can be run in two modes - in the mode of a single query and in the mode of reading the input stream

### Single mode

Server address, password and command to server must be specified in flags at startup. Example:

    ./rcon -a 127.0.0.1:16260 -p mypassword -c help
    
If flags are passed into the CLI, then the flags are sent as a single command. 
The response is displayed, and the CLI will exit.

### Input stream mode

To run CLI in input stream mode add command cli. Example:

    ./rcon -a 127.0.0.1:16260 -p mypassword cli
    
Use `^C` to terminate or type command `:q` in CLI to exit.    

## Configuration file

For more convenient use, the ability to create the rcon.yaml configuration file is provided. 
You can save the address of the remote server and its password. If the configuration file exists, 
and in it the default block is filled, then at startup the -a and -p flags can be omitted. Examples:

    ./rcon -a 127.0.0.1:16260 -c players
    ./rcon -c status
    ./rcon -p mypassword cli
    ./rcon cli

Default configuration file name is `rcon.yaml`. File must be saved in yaml format. It is also possible 
to set the environment name and connection parameters for each server, for example:

```yaml
default:
  address: "127.0.0.1:16260"
  password: "password"
zomboid:
  address: "127.0.0.1:16260"
  password: "password"
rust:
  address: "127.0.0.1:28003"
  password: "password"  
```

You can choose the environment at the start:

    ./rcon -e rust -c status
    ./rcon -e zomboid cli
