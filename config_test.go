package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/gorcon/rcon-cli/internal/session"
	"github.com/stretchr/testify/assert"
)

// DefaultTestLogName sets the default log file name.
const DefaultTestLogName = "rcon-test.log"

const ConfigLayoutJSON = `{"%s": {"address": "%s", "password": "%s", "log": "%s", "type": "%s"}}`
const ConfigLayoutYAML = "%s:\n  address: %s\n  password: %s\n  log: %s\n  type: %s"

func TestNewConfig(t *testing.T) {
	t.Run("no errors yaml", func(t *testing.T) {
		expected := Config{
			"default": session.Session{Address: "", Password: "", Log: "rcon-default.log"},
		}

		cfg, err := NewConfig("rcon.yaml")
		assert.NoError(t, err)
		assert.Equal(t, &expected, cfg)
	})

	t.Run("no errors json", func(t *testing.T) {
		configFileName := "rcon-test-local.json"
		stringBody := fmt.Sprintf(ConfigLayoutJSON, DefaultConfigEnv, "", "", DefaultTestLogName, "")
		err := createFile(configFileName, stringBody)
		assert.NoError(t, err)

		defer func() {
			err := os.Remove(configFileName)
			assert.NoError(t, err)
		}()

		expected := Config{
			DefaultConfigEnv: session.Session{Address: "", Password: "", Log: DefaultTestLogName},
		}

		cfg, err := NewConfig(configFileName)
		assert.NoError(t, err)
		assert.Equal(t, &expected, cfg)
	})

	t.Run("file not exists", func(t *testing.T) {
		cfg, err := NewConfig("nonexist.yaml")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("unexpected error: %v", err)
		}

		assert.Nil(t, cfg)
	})

	// Test is valid because of automatic placement of a temporary binary to the
	// `/tmp` directory.
	// Expected error message: `read config error: open /tmp/rcon.yaml: no such file or directory`.
	t.Run("default file not exists", func(t *testing.T) {
		cfg, err := NewConfig("")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("unexpected error: %v", err)
		}

		assert.Nil(t, cfg)
	})

	t.Run("file is incorrect", func(t *testing.T) {
		configFileName := "rcon-test-local.yaml"
		stringBody := fmt.Sprintf("address: \"%s\"\n  password: \"%s\"\n  log: \"%s\"", "", MockPasswordRCON, DefaultTestLogName)
		err := createFile(configFileName, stringBody)
		assert.NoError(t, err)

		defer func() {
			err := os.Remove(configFileName)
			assert.NoError(t, err)
		}()

		cfg, err := NewConfig(configFileName)
		assert.EqualError(t, err, "read config error: yaml: line 1: did not find expected key")

		assert.Nil(t, cfg)
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
