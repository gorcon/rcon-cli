# rcon-cli
[![Go Report Card](https://goreportcard.com/badge/github.com/gorcon/rcon-cli)](https://goreportcard.com/report/github.com/gorcon/rcon-cli)

CLI for executing queries on a remote server

## Supported Games

* [Project Zomboid](https://store.steampowered.com/app/108600) 
* [Conan Exiles](https://store.steampowered.com/app/440900)
* [Rust](https://store.steampowered.com/app/252490) (add +rcon.web 0 to the args when starting the server)
* [ARK: Survival Evolved](https://store.steampowered.com/app/346110)

Open pull request if you have successfully used a package with another game with rcon support and add it to the list.

## Installation

Download the binary for your platform from the [releases](https://github.com/gorcon/rcon/releases)

## Usage

Rcon CLI can be run in two modes - in the mode of a single query and in the mode of reading the input stream

### Single mode

Server address, password and command to server must be specified in flags at startup. Example:

    ./rcon -a address -p password -c command
    
If flags are passed into the CLI, then the flags are sent as a single command. 
The response is displayed, and the CLI will exit.

### Input stream mode

To run CLI in input stream mode add command cli. Example:

    ./rcon -a address -p password cli
    
Use ^C to terminate or type command :q in CLI.    

## Configuration file

For more convenient use, the ability to create the rcon.yaml configuration file is provided. 
You can save the address of the remote server and its password. If the configuration file exists, 
then at startup the -a and -p flags can be omitted. 

File must be saved in yaml format. Example:

```yaml
default:
  address: "127.0.0.1:16260"
  password: "password"
```
