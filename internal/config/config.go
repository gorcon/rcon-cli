package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DefaultConfigName sets the default config file name.
const DefaultConfigName = "rcon.yaml"

// DefaultConfigEnv is the name of the environment, which is taken
// as default unless another value is passed.
const DefaultConfigEnv = "default"

var (
	// ErrConfigValidation is when config validation completed with errors.
	ErrConfigValidation = errors.New("config validation error")

	// ErrUnsupportedFileExt is returned when config file has an unsupported
	// extension. Allowed extensions is `.json`, `.yml`, `.yaml`.
	ErrUnsupportedFileExt = errors.New("unsupported file extension")
)

// Config allows to take a remote server address and password from
// the configuration file. This enables not to specify these flags when
// running the CLI.
//
// Example:
// ```yaml
// default:
//
//	address: "127.0.0.1:16260"
//	password: "password"
//
// ```.
type Config map[string]Session

// NewConfig finds and parses config file with remote server credentials.
func NewConfig(name string) (*Config, error) {
	cfg := new(Config)
	if err := cfg.ParseFromFile(name); err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// ParseFromFile reads a configuration file from disk and loads its contents into
// the application's config structure. YAML and JSON files are supported.
func (cfg *Config) ParseFromFile(name string) error {
	if name != "" {
		return cfg.parse(name)
	}

	home, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fmt.Errorf("get abs path: %w", err)
	}

	name = home + "/" + DefaultConfigName
	if err = cfg.parse(name); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	*cfg = Config{DefaultConfigEnv: {}}

	return nil
}

// Validate validates the config fields.
func (cfg *Config) Validate() error {
	if cfg == nil {
		return fmt.Errorf("%w: config is not set", ErrConfigValidation)
	}

	for key, ses := range *cfg {
		switch ses.Type {
		case "", ProtocolRCON, ProtocolTELNET, ProtocolWebRCON:
		default:
			return fmt.Errorf("%w: unsupported type in %s environment", ErrConfigValidation, key)
		}
	}

	return nil
}

func (cfg *Config) parse(name string) error {
	file, err := os.ReadFile(name)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	switch ext := path.Ext(name); ext {
	case ".yml", ".yaml":
		err = yaml.Unmarshal(file, cfg)
	case ".json":
		err = json.Unmarshal(file, cfg)
	default:
		err = fmt.Errorf("%w %s", ErrUnsupportedFileExt, ext)
	}

	return err
}
