package logger

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// LogRecordTimeLayout is layout for convert time.Now to String.
const LogRecordTimeLayout = "2006-01-02 15:04:05"

// LogRecordFormat is format to log line record.
const LogRecordFormat = "[%s] %s: %s\n%s\n\n"

// GetLogFile opens file for append strings. Creates file if file not exist.
func GetLogFile(logName string) (*os.File, error) {
	if logName == "" {
		return nil, errors.New("empty file name")
	}

	var file *os.File

	_, err := os.Stat(logName)

	switch {
	case err == nil:
		// Open current file.
		file, err = os.OpenFile(logName, os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			return nil, err
		}
	case os.IsNotExist(err):
		// Create new file.
		file, err = os.Create(logName)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	return file, nil
}

// AddLog saves request and response to log file.
func AddLog(logName string, address string, request string, response string) error {
	// Disable logging if log file name is empty.
	if logName == "" {
		return nil
	}

	file, err := GetLogFile(logName)
	if err != nil {
		return err
	}
	defer file.Close()

	line := fmt.Sprintf(LogRecordFormat, time.Now().Format(LogRecordTimeLayout), address, request, response)
	if _, err := file.WriteString(line); err != nil {
		return err
	}

	return nil
}
