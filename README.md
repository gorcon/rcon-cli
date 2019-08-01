# rcon-cli
[![Go Report Card](https://goreportcard.com/badge/github.com/gorcon/rcon-cli)](https://goreportcard.com/report/github.com/gorcon/rcon-cli)
[![GitHub All Releases](https://img.shields.io/github/downloads/gorcon/rcon-cli/total)](https://github.com/gorcon/rcon-cli/releases)

CLI for executing queries on a remote server

## Supported Games

* [Project Zomboid](https://store.steampowered.com/app/108600) 
* [Conan Exiles](https://store.steampowered.com/app/440900)
* [Rust](https://store.steampowered.com/app/252490) (add +rcon.web 0 to the args when starting the server)
* [ARK: Survival Evolved](https://store.steampowered.com/app/346110)

Open pull request if you have successfully used a package with another game with rcon support and add it to the list.

## Installation

Download the binary for your platform from the [latest releases](https://github.com/gorcon/rcon-cli/releases/latest)

See [Changelog](changelog.md) for release details

## Usage

```text
USAGE:
   rcon [global options] command [command options] [arguments...]

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -a value, --address value  set host and port to remote rcon server. Example 127.0.0.1:16260
                              can be set in the config file rcon.yaml
   -p value, --password value  set password to remote rcon server
                               can be set in the config file rcon.yaml
   -c value, --command value  command to execute on remote server. Required flag to run in single mode
   -e value, --env value      allows to select remote server address and password from the environment
                              in the configuration file
   -l value, --log value  path and name of the log file. if not specified, it is taken from the config.
   --cfg value            allows to specify the path and name of the configuration file. The default
                value is rcon.yaml.
   --help, -h     show help
   --version, -v  print the version
```

Rcon CLI can be run in two modes - in the mode of a single query and in the mode of reading the input stream

### Single mode

Server address, password and command to server must be specified in flags at startup. Example:

    ./rcon -a 127.0.0.1:16260 -p mypassword -c help
    
If arguments are passed they will sent as a single command. The response will displayed, and the CLI will exit.

### Interactive input stream mode

To run CLI in interactive mode run `rcon` without `-c` argument. Example:

    ./rcon -a 127.0.0.1:16260 -p mypassword
    
Use `^C` to terminate or type command `:q` to exit.    

## Configuration file

For more convenient use, the ability to create the rcon.yaml configuration file is provided. 
You can save the host and port of the remote server and its password. If the configuration file exists, 
and in it the default block is filled, then at startup the -a and -p flags can be omitted. Examples:

    ./rcon -a 127.0.0.1:16260 -c players
    ./rcon -c status
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
```

You can choose the environment at the start:

    ./rcon -e rust -c status
    ./rcon -e zomboid
    
And set custom config file     

    ./rcon --cfg /path/to/config/file.yaml
    
You can use `-l` argument to specify path to log file.

     ./rcon -l /path/to/file.log

## Contribute

If you think that you have found a bug, create an issue and indicate your operating system, platform and the game on 
which the error was reproduced. Also describe what you were doing so that the error could be reproduced.

## License

MIT License, see [LICENSE](LICENSE)
