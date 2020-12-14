package config_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/gorcon/rcon-cli/internal/config"
	"github.com/gorcon/rcon-cli/internal/session"
	"github.com/stretchr/testify/assert"
)

// DefaultTestLogName sets the default log file name.
const DefaultTestLogName = "rcon-test.log"

const ConfigLayoutJSON = `{"%s": {"address": "%s", "password": "%s", "log": "%s", "type": "%s"}}`
const ConfigLayoutYAML = "%s:\n  address: %s\n  password: %s\n  log: %s\n  type: %s"

func TestNewConfig(t *testing.T) {
	t.Run("no errors yaml", func(t *testing.T) {
		configFileName := "rcon-test-local.yaml"
		stringBody := fmt.Sprintf(ConfigLayoutYAML, config.DefaultConfigEnv, "", "", DefaultTestLogName, "")
		createFile(configFileName, stringBody)
		defer os.Remove(configFileName)

		expected := config.Config{
			"default": session.Session{Address: "", Password: "", Log: "rcon-test.log"},
		}

		cfg, err := config.NewConfig(configFileName)
		assert.NoError(t, err)
		assert.Equal(t, &expected, cfg)
	})

	t.Run("no errors json", func(t *testing.T) {
		configFileName := "rcon-test-local.json"
		stringBody := fmt.Sprintf(ConfigLayoutJSON, config.DefaultConfigEnv, "", "", DefaultTestLogName, "")
		createFile(configFileName, stringBody)
		defer os.Remove(configFileName)

		expected := config.Config{
			config.DefaultConfigEnv: session.Session{Address: "", Password: "", Log: DefaultTestLogName},
		}

		cfg, err := config.NewConfig(configFileName)
		assert.NoError(t, err)
		assert.Equal(t, &expected, cfg)
	})

	t.Run("file not exists", func(t *testing.T) {
		cfg, err := config.NewConfig("nonexist.yaml")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("unexpected error: %v", err)
		}

		assert.Nil(t, cfg)
	})

	// Test is valid because of automatic placement of a temporary binary to the
	// `/tmp` directory.
	// Expected error message: `read config error: open /tmp/rcon.yaml: no such file or directory`.
	t.Run("default file not exists", func(t *testing.T) {
		cfg, err := config.NewConfig("")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("unexpected error: %v", err)
		}

		assert.Nil(t, cfg)
	})

	t.Run("file is incorrect", func(t *testing.T) {
		configFileName := "rcon-test-local.yaml"
		stringBody := fmt.Sprintf("address: \"%s\"\n  password: \"%s\"\n  log: \"%s\"", "", "password", DefaultTestLogName)
		createFile(configFileName, stringBody)
		defer os.Remove(configFileName)

		cfg, err := config.NewConfig(configFileName)
		assert.EqualError(t, err, "read config error: yaml: line 1: did not find expected key")

		assert.Nil(t, cfg)
	})

	t.Run("unsupported file extension", func(t *testing.T) {
		configFileName := "unsupported-local.ini"
		stringBody := "[genera]\addr="
		createFile(configFileName, stringBody)
		defer os.Remove(configFileName)

		cfg, err := config.NewConfig(configFileName)
		assert.EqualError(t, err, "read config error: unsupported file extension .ini")

		assert.Nil(t, cfg)
	})

	t.Run("validation failed", func(t *testing.T) {
		configFileName := "rcon-test-local.json"
		stringBody := fmt.Sprintf(ConfigLayoutJSON, config.DefaultConfigEnv, "", "", DefaultTestLogName, "pigeon post")
		createFile(configFileName, stringBody)
		defer os.Remove(configFileName)

		cfg, err := config.NewConfig(configFileName)
		assert.EqualError(t, err, "config validation error: unsupported type in default environment")

		expected := config.Config{
			config.DefaultConfigEnv: session.Session{Address: "", Password: "", Log: DefaultTestLogName, Type: "pigeon post"},
		}

		assert.Equal(t, &expected, cfg)
	})
}

func TestConfig_Validate(t *testing.T) {
	t.Run("initialized empty config", func(t *testing.T) {
		cfg := new(config.Config)
		err := cfg.Validate()
		assert.NoError(t, err)
	})

	t.Run("not initialized empty config", func(t *testing.T) {
		var cfg *config.Config
		err := cfg.Validate()
		assert.EqualError(t, err, "config validation error: config is not set")
	})
}

func createFile(name, stringBody string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}

	_, err = file.WriteString(stringBody)

	return err
}
