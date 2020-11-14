package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gorcon/rcon-cli/internal/session"
	"gopkg.in/yaml.v3"
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
// ```.
type Config map[string]session.Session

// GetConfig finds and parses config file for details of connecting to
// a remote server.
func GetConfig(path string) (*Config, error) {
	if path == "" {
		home, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return nil, err
		}

		path = home + "/" + DefaultConfigName
		if _, err := os.Stat(path); err != nil {
			return nil, err
		}
	}

	// Read the config file if file exists.
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	cfg, err := ReadYamlConfig(path)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

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
