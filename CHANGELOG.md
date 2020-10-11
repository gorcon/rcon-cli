# Changelog
All notable changes to this project will be documented in this file.

**ATTN**: This project uses [semantic versioning](http://semver.org/).

## [Unreleased]
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

[Unreleased]: https://github.com/gorcon/rcon-cli/compare/v0.7.0...HEAD
[v0.7.0]: https://github.com/gorcon/rcon-cli/compare/0.6.0...v0.7.0
[v0.6.0]: https://github.com/gorcon/rcon-cli/compare/0.5.0...v0.6.0
[v0.5.0]: https://github.com/gorcon/rcon-cli/compare/v0.4.0...0.5.0
[v0.4.0]: https://github.com/gorcon/rcon-cli/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/gorcon/rcon-cli/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/gorcon/rcon-cli/compare/v0.1.0...v0.2.0
