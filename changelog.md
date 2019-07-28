# Changelog
All notable changes to this project will be documented in this file.


## [v0.3.0] - 2019-07-28
### Added
- Print error messages for missed address and password arguments.
- Added argument `--cfg` to pass custom config path/name.
- Check of parameters of authorization of a remote server before launching interactive mode. 
- Added rcon.yaml as config sample.
- Added log variable to config file that enables logs for requests and responses for remote server.

### Removed
- Remove rcon-upx fron release.
- Remove 'cli' command to run in interactive mode. For use interactive mode run `rcon` without `-c` argument.


## [v0.2.0] - 2019-07-27
### Added
- Added environments to config.

### Fixed
- Fix global options in interactive mode.
