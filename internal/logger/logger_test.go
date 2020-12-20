package logger_test

import (
	"os"
	"testing"

	"github.com/gorcon/rcon-cli/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestOpenFile(t *testing.T) {
	logDir := "temp"
	logName := "tmpfile.log"
	logPath := logDir + "/" + logName

	// Test empty log file name.
	t.Run("empty file name", func(t *testing.T) {
		file, err := logger.OpenFile("")
		assert.Nil(t, file)
		assert.EqualError(t, err, "empty file name")
	})

	// Positive test create new log file.
	t.Run("create new log file", func(t *testing.T) {
		os.Mkdir(logDir, 0700)
		defer os.RemoveAll(logDir)

		file, err := logger.OpenFile(logPath)
		assert.NotNil(t, file)
		assert.NoError(t, err)
	})
}

func TestWrite(t *testing.T) {
	logName := "tmpfile.log"

	address := "127.0.0.1:16200"
	command := "players"
	result := `Players connected (2):
-admin
-testuser`

	defer os.Remove(logName)

	// Test skip log. No logs is available.
	t.Run("skip log", func(t *testing.T) {
		err := logger.Write("", address, command, result)
		assert.NoError(t, err)
	})

	// Test create log file.
	t.Run("create log file", func(t *testing.T) {
		err := logger.Write(logName, address, command, result)
		assert.NoError(t, err)
	})

	// Test append to log file.
	t.Run("append to log file", func(t *testing.T) {
		err := logger.Write(logName, address, command, result)
		assert.NoError(t, err)
	})
}
