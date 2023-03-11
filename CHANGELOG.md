# Changelog
All notable changes to this project will be documented in this file.

**ATTN**: This project uses [semantic versioning](http://semver.org/).

## [Unreleased]

## [v0.10.3] - 2023-03-11
## Added
- Added folder creation for logs.

### Fixed
- Fixed default value to config flag.

## Updated
- Updated Go modules (go1.19).
- Updated golang-ci linter (1.50.1).
- Updated dependencies.

## [v0.10.2] - 2022-03-09
### Added
- Added `--timeout, -T` flag, allowed to set dial and execute timeout.
- Added Makefile.

## [v0.10.1] - 2021-11-13
### Fixed
- Fixed close connection panic #19 

## [v0.10.0] - 2021-09-09
### Added
- Added `--skip, -s` flag, allowed to skip error on multiple commands or in terminal mode.

### Changed
- Create one keep alive connection in terminal mode.

## [v0.9.1] - 2021-03-14
### Added
- Added Dockerfile. 

### Changed
- Disabled CGO in release script.

## [v0.9.0] - 2020-12-20
### Added
- Added config validation. 
- Added config JSON format supporting. 
- Added protocol type asking in interactive mode.
- Added protocol type validation in interactive mode.
- Added the ability to send several commands in a row with one request #10

### Changed
- Removed `-c value, --command value` flag. It replaced to commands without flags.
- Flag `--cfg value` replaced to `-c value, --config value` flag. 

## [v0.8.1] - 2020-11-17
### Added
- Added tests for real servers. Servers list: `Project Zomboid`, `7 Days to Die`, `Rust`. 
- Added removing part of the constantly repeated data from `7 Days to Die` response ([details](https://github.com/gorcon/telnet/issues/1)). 

## [v0.8.0-beta.2] - 2020-10-18
### Added
- Added interactive mode for Web RCON #12.

### Fixed
- Fixed response for another request for Rust server #13.

## [v0.8.0-beta] - 2020-10-18
### Added
- Added Rust Web RCON support #8. Add `-t web` argument when execute `rcon` cli.

### Changed
- Code and tests refactoring.

### Fixed
- Fixed changelog.

## [v0.7.0] - 2020-10-10
### Added
- Added support amd64 darwin compilation.
- Added 7 Days to Die support #5. Add `-t telnet` argument when execute `rcon` cli.

## [v0.6.0] - 2020-07-10
### Fixed
- Fix `Conan Exiles "response for another request" error` bug #6.

## [v0.5.0] - 2020-06-18
### Fixed
- Correction of text after AutoCorrect.

### Added
- More tests.
- Add Go modules (go1.13).
- Add golang-ci linter.

## [v0.4.0] - 2019-08-05
### Added
- Added argument `-l`, `--log` to pass custom log path/name. Argument has higher priority .
than entry in configuration file.
- Added interactive opportunity to enter server address and password.

## [v0.3.0] - 2019-07-28
### Added
- Print error messages for missed address and password arguments.
- Added argument `--cfg` to pass custom config path/name.
- Check of parameters of authorization of a remote server before launching interactive mode. 
- Added rcon.yaml as config sample.
- Added log variable to config file that enables logs for requests and responses for remote server.

### Removed
- Remove `rcon-upx` from release.
- Remove 'cli' command to run in interactive mode. For use interactive mode run `rcon` without `-c` argument.

## [v0.2.0] - 2019-07-27
### Added
- Added environments to config.

### Fixed
- Fix global options in interactive mode.

## v0.1.0 - 2019-07-22
### Added
- Initial implementation.

[Unreleased]: https://github.com/gorcon/rcon-cli/compare/v0.10.3...HEAD
[v0.10.3]: https://github.com/gorcon/rcon-cli/compare/v0.10.2...v0.10.3
[v0.10.2]: https://github.com/gorcon/rcon-cli/compare/v0.10.1...v0.10.2
[v0.10.1]: https://github.com/gorcon/rcon-cli/compare/v0.10.0...v0.10.1
[v0.10.0]: https://github.com/gorcon/rcon-cli/compare/v0.9.1...v0.10.0
[v0.9.1]: https://github.com/gorcon/rcon-cli/compare/v0.9.0...v0.9.1
[v0.9.0]: https://github.com/gorcon/rcon-cli/compare/v0.8.1...v0.9.0
[v0.8.1]: https://github.com/gorcon/rcon-cli/compare/v0.8.0-beta.2...v0.8.1
[v0.8.0-beta.2]: https://github.com/gorcon/rcon-cli/compare/v0.8.0-beta...v0.8.0-beta.2
[v0.8.0-beta]: https://github.com/gorcon/rcon-cli/compare/v0.7.0...v0.8.0-beta
[v0.7.0]: https://github.com/gorcon/rcon-cli/compare/v0.6.0...v0.7.0
[v0.6.0]: https://github.com/gorcon/rcon-cli/compare/0.5.0...v0.6.0
[v0.5.0]: https://github.com/gorcon/rcon-cli/compare/v0.4.0...0.5.0
[v0.4.0]: https://github.com/gorcon/rcon-cli/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/gorcon/rcon-cli/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/gorcon/rcon-cli/compare/v0.1.0...v0.2.0
