package logger

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// DefaultTimeLayout is layout for convert time.Now to String.
const DefaultTimeLayout = "2006-01-02 15:04:05"

// DefaultLineFormat is format to log line record.
const DefaultLineFormat = "[%s] %s: %s\n%s\n\n"

// ErrEmptyFileName is returned when trying to open file with empty name.
var ErrEmptyFileName = errors.New("empty file name")

// OpenFile opens file for append strings. Creates file if file not exist.
func OpenFile(name string) (file *os.File, err error) {
	if name == "" {
		return nil, ErrEmptyFileName
	}

	switch _, err = os.Stat(name); {
	case err == nil:
		file, err = os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return file, fmt.Errorf("open: %w", err)
		}
	case os.IsNotExist(err):
		file, err = os.Create(name)
		if err != nil {
			return file, fmt.Errorf("create: %w", err)
		}
	}

	return file, nil
}

// Write saves request and response to log file.
func Write(name string, address string, request string, response string) error {
	// Disable logging if log file name is empty.
	if name == "" {
		return nil
	}

	file, err := OpenFile(name)
	if err != nil {
		return err
	}
	defer file.Close()

	line := fmt.Sprintf(DefaultLineFormat, time.Now().Format(DefaultTimeLayout), address, request, response)
	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
