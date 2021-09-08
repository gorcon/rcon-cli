# rcon-cli
[![GitHub Build](https://github.com/gorcon/rcon-cli/workflows/build/badge.svg)](https://github.com/gorcon/rcon-cli/actions)
[![top level coverage](https://gocover.io/_badge/github.com/gorcon/rcon-cli?0)](https://gocover.io/github.com/gorcon/rcon-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/gorcon/rcon-cli)](https://goreportcard.com/report/github.com/gorcon/rcon-cli)
[![GitHub All Releases](https://img.shields.io/github/downloads/gorcon/rcon-cli/total)](https://github.com/gorcon/rcon-cli/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/outdead/rcon.svg)](https://hub.docker.com/r/outdead/rcon)

CLI for executing queries on a remote [Source dedicated game server](https://developer.valvesoftware.com/wiki/Source_Dedicated_Server), using the [RCON](https://developer.valvesoftware.com/wiki/Source_RCON_Protocol) protocol.

## Supported Games

* [Project Zomboid](https://store.steampowered.com/app/108600) 
* [Conan Exiles](https://store.steampowered.com/app/440900)
* [Rust](https://store.steampowered.com/app/252490) (add `+rcon.web 0` to the args when starting the server or add `-t web` to `rcon-cli` args)
* [ARK: Survival Evolved](https://store.steampowered.com/app/346110)
* [7 Days to Die](https://store.steampowered.com/app/251570) (add `-t telnet` to rcon-cli args)

Open pull request if you have successfully used a package with another game with rcon support and add it to the list.

## Installation

Download the binary for your platform from the [latest releases](https://github.com/gorcon/rcon-cli/releases/latest)

See [Changelog](CHANGELOG.md) for release details

### Docker

    docker pull outdead/rcon

## Usage

```text
USAGE:
   rcon [options] [commands...]

GLOBAL OPTIONS:
   --address value, -a value   Set host and port to remote server. Example 127.0.0.1:16260
   --password value, -p value  Set password to remote server
   --type value, -t value      Specify type of connection (default: rcon)
   --log value, -l value       Path to the log file. If not specified it is taken from the config
   --config value, -c value    Path to the configuration file (default: rcon.yaml)
   --env value, -e value       Config environment with server credentials (default: default)
   --skip, -s                  Skip errors and run next command (default: false)   
   --help, -h                  show help (default: false)
   --version, -v               print the version (default: false)
```

Rcon CLI can be run in two modes - in the mode of a single query and in the mode of reading the input stream

### Single mode

Server address, password and command to server must be specified in flags at startup. Example:  

    ./rcon -a 127.0.0.1:16260 -p mypassword command
    
Since from `rcon-cli` 0.9.0 version it is possible to send several commands in one request. Example:  

    ./rcon -a 127.0.0.1:16260 -p mypassword command "command with several words" 'command "with double quotes"'

If commands are passed they will sent in a single mode. The response will displayed, and the CLI will exit.

### Interactive input stream mode

To run CLI in interactive mode run `rcon` without commands. Example:

    ./rcon -a 127.0.0.1:16260 -p mypassword
    
Use `^C` to terminate or type command `:q` to exit.    

## Configuration file

For more convenient use, the ability to create the rcon.yaml configuration file is provided. 
You can save the host and port of the remote server and its password. If the configuration file exists, 
and in it the default block is filled, then at startup the -a and -p flags can be omitted. Examples:

    ./rcon -a 127.0.0.1:16260 players
    ./rcon status
    ./rcon -p mypassword
    ./rcon 

Default configuration file name is `rcon.yaml`. File must be saved in yaml format. It is also possible 
to set the environment name and connection parameters for each server. You can enable logging requests 
and responses. To do this, you need to define the log variable in the environment blocks. You can do 
this for each server separately and create different log files for them. If the path to the log file is 
not specified, then logging will not be conducted. 

```yaml
default:
  address: "127.0.0.1:16260"
  password: "password"
  log: "rcon-default.log"
zomboid:
  address: "127.0.0.1:16260"
  password: "password"
  log: "rcon-zomboid.log"
rust:
  address: "127.0.0.1:28003"
  password: "password"
7dtd:
  address: "172.19.0.2:8081"
  password: "password"
  type: "telnet"
```

You can choose the environment at the start:

    ./rcon -e rust status
    ./rcon -e zomboid
    
And set custom config file     

    ./rcon -c /path/to/config/file.yaml
    
You can use `-l` argument to specify path to log file.

     ./rcon -l /path/to/file.log

Since from `rcon-cli` 0.7.0 version support for the TELNET protocol has been added. On this protocol remote access to 
the 7 Days to Die console is based. You can use `-t telnet` argument to specify the protocol type.

     ./rcon -a 172.19.0.2:8081 -p password -t telnet version
     
Since from `rcon-cli` 0.8.0 version support for the Web RCON protocol has been added. On this protocol remote access to 
the Rust console is based. You can use `-t web` argument to specify the protocol type.

     ./rcon -a 127.0.0.1:28016 -p password -t web status

## Contribute

If you think that you have found a bug, create an issue and indicate your operating system, platform and the game on 
which the error was reproduced. Also describe what you were doing so that the error could be reproduced.

## License

MIT License, see [LICENSE](LICENSE)
