# rcon-cli
[![Mentioned in Awesome-Selfhosted](https://awesome.re/mentioned-badge.svg)](https://github.com/awesome-selfhosted/awesome-selfhosted)
[![GitHub Build](https://github.com/gorcon/rcon-cli/workflows/build/badge.svg)](https://github.com/gorcon/rcon-cli/actions)
[![top level coverage](https://gocover.io/_badge/github.com/gorcon/rcon-cli?0)](https://gocover.io/github.com/gorcon/rcon-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/gorcon/rcon-cli)](https://goreportcard.com/report/github.com/gorcon/rcon-cli)
[![GitHub All Releases](https://img.shields.io/github/downloads/gorcon/rcon-cli/total)](https://github.com/gorcon/rcon-cli/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/outdead/rcon.svg)](https://hub.docker.com/r/outdead/rcon)

CLI for executing queries on a remote [Source dedicated game server](https://developer.valvesoftware.com/wiki/Source_Dedicated_Server), using the [RCON](https://developer.valvesoftware.com/wiki/Source_RCON_Protocol) protocol.

## Supported Games
* [7 Days to Die](https://store.steampowered.com/app/251570) (add `-t telnet` to rcon-cli args)
* [ARK: Survival Evolved](https://store.steampowered.com/app/346110)
* [Conan Exiles](https://store.steampowered.com/app/440900)
* [Counter-Strike: Global Offensive](https://store.steampowered.com/app/730)
* [Factorio](https://factorio.com/)
* [Minecraft](https://www.minecraft.net)
* [Project Zomboid](https://store.steampowered.com/app/108600) 
* [Rust](https://store.steampowered.com/app/252490) (add `+rcon.web 0` to the args when starting the server or add `-t web` to `rcon-cli` args)
* [Team Fortress 2](https://store.steampowered.com/app/440/Team_Fortress_2/)
* [V Rising](https://store.steampowered.com/app/1604030/V_Rising/)




Open pull request if you have successfully used a package with another game with rcon support and add it to the list.

## Installation
Download the binary for your platform from the [latest releases](https://github.com/gorcon/rcon-cli/releases/latest)

See [Changelog](CHANGELOG.md) for release details

### Docker
```bash
docker pull outdead/rcon
```

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
   --timeout value, -T value   Set dial and execute timeout (default: 10s)
   --help, -h                  show help (default: false)
   --version, -v               print the version (default: false)
```

Rcon CLI can be run in two modes - in the mode of a single query and in the mode of reading the input stream

### Single mode
Server address, password and command to server must be specified in flags at startup. Example:  
```bash
./rcon -a 127.0.0.1:16260 -p mypassword command
```

It is possible to send several commands in one request. Example:  
```bash
./rcon -a 127.0.0.1:16260 -p mypassword command "command with several words" 'command "with double quotes"'
```

If commands passed, they sent in a single mode. The response displayed, and the CLI will exit.

### Interactive input stream mode
To run CLI in interactive mode run `rcon` without commands. Example:
```bash
./rcon -a 127.0.0.1:16260 -p mypassword
```

Use `^C` to terminate or type command `:q` to exit.    

### In Docker
```bash
docker run -it --rm outdead/rcon ./rcon [options] [commands...]
```

You can add your config file as volume:
```bash
docker run -it --rm \
      -v /path/to/rcon-local.yaml:/rcon.yaml \
      outdead/rcon ./rcon -c rcon.yaml -e default players
```

## Configuration file
For more convenient use, the ability to create the `rcon.yaml` configuration file provided. You can save the host and port of the remote server and its password. If the configuration file exists, and the default block filled in it, then at startup the `-a` and `-p` flags can be omitted. Examples:
```bash
./rcon -a 127.0.0.1:16260 players
./rcon status
./rcon -p mypassword
./rcon
```

Default configuration file name is `rcon.yaml`. File must be saved in yaml format. It is also possible to set the environment name and connection parameters for each server. You can enable logging requests and responses. To do this, you need to define the log variable in the environment blocks. You can do 
this for each server separately and create different log files for them. If the path to the log file not specified, then logging will not be conducted. 
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

## Args
You can choose the environment at the start:
```bash
./rcon -e rust status
./rcon -e zomboid
```

Set custom config file:
```bash
./rcon -c /path/to/config/file.yaml
```

Use `-l` argument to specify path to log file:
```bash
./rcon -l /path/to/file.log
```

Use `-t` argument to specify the protocol type:
```bash
# 7 Days to Die
./rcon -a 172.19.0.2:8081 -p password -t telnet version

# Rust
./rcon -a 127.0.0.1:28016 -p password -t web status
```

Use `-T` argument to specify dial and execute timeout:
```bash
./rcon -a 172.19.0.2:8081 -p password -t telnet -T 10s version
```

## Contribute
If you think that you have found a bug, create an issue and indicate your operating system, platform, and the game on which the error reproduced. Also describe what you were doing so that the error could be reproduced.

## License
MIT License, see [LICENSE](LICENSE)
