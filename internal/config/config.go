package config

import (
	"io/ioutil"

	"github.com/gorcon/rcon-cli/internal/session"
	"gopkg.in/yaml.v2"
)

// Config allows to take a remote server address and password from
// the configuration file. This enables not to specify these flags when
// running the CLI.
//
// Example:
// ```yaml
// default:
//   address: "127.0.0.1:16260"
//   password: "password"
// ```.
type Config map[string]session.Session

// ReadYamlConfig reads config data from yaml file.
func ReadYamlConfig(path string) (cfg Config, err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
