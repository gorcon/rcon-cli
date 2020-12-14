package logger_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gorcon/rcon-cli/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestAddLog(t *testing.T) {
	logName := "tmpfile.log"

	address := "127.0.0.1:16200"
	command := "players"
	result := `Players connected (2):
-admin
-testuser`

	defer func() {
		err := os.Remove(logName)
		assert.NoError(t, err)
	}()

	// Test skip log. No logs is available.
	t.Run("skip log", func(t *testing.T) {
		err := logger.AddLog("", address, command, result)
		assert.NoError(t, err)
	})

	// Test create log file.
	t.Run("create log file", func(t *testing.T) {
		err := logger.AddLog(logName, address, command, result)
		assert.NoError(t, err)
	})

	// Test append to log file.
	t.Run("append to log file", func(t *testing.T) {
		err := logger.AddLog(logName, address, command, result)
		assert.NoError(t, err)
	})
}

func TestGetLogFile(t *testing.T) {
	logDir := "temp"
	logName := "tmpfile.log"
	logPath := logDir + "/" + logName

	// Test empty log file name.
	t.Run("empty file name", func(t *testing.T) {
		file, err := logger.GetLogFile("")
		assert.Nil(t, file)
		assert.EqualError(t, err, "empty file name")
	})

	// Test stat permission denied.
	t.Run("stat permission denied", func(t *testing.T) {
		if err := os.Mkdir(logDir, 0400); err != nil {
			assert.NoError(t, err)
			return
		}
		defer func() {
			err := os.RemoveAll(logDir)
			assert.NoError(t, err)
		}()

		file, err := logger.GetLogFile(logPath)
		assert.Nil(t, file)
		assert.EqualError(t, err, fmt.Sprintf("stat %s: permission denied", logPath))
	})

	// Test create permission denied.
	t.Run("open permission denied", func(t *testing.T) {
		if err := os.Mkdir(logDir, 0500); err != nil {
			assert.NoError(t, err)
			return
		}
		defer func() {
			err := os.RemoveAll(logDir)
			assert.NoError(t, err)
		}()

		file, err := logger.GetLogFile(logPath)
		assert.Nil(t, file)
		assert.EqualError(t, err, fmt.Sprintf("open %s: permission denied", logPath))
	})

	// Positive test create new log file + test open permission denied.
	t.Run("create new log file", func(t *testing.T) {
		if err := os.Mkdir(logDir, 0700); err != nil {
			assert.NoError(t, err)
			return
		}
		defer func() {
			err := os.RemoveAll(logDir)
			assert.NoError(t, err)
		}()

		// Positive test create new log file.
		file, err := logger.GetLogFile(logPath)
		assert.NotNil(t, file)
		assert.NoError(t, err)

		if err := os.Chmod(logPath, 0000); err != nil {
			assert.NoError(t, err)
			return
		}

		// Test open permission denied.
		file, err = logger.GetLogFile(logPath)
		assert.Nil(t, file)
		assert.EqualError(t, err, fmt.Sprintf("open %s: permission denied", logPath))
	})
}
