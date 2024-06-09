package log

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type CustomLogger struct {
	Logger *log.Logger
	file   *os.File
}

type Message struct {
	Level        string `json:"level,omitempty"`
	TraceId      string `json:"traceId,omitempty"`
	Method       string `json:"method,omitempty"`
	URI          string `json:"uri,omitempty"`
	Msg          string `json:"msg,omitempty"`
	StatusCode   int    `json:"statusCode,omitempty"`
	Duration     int64  `json:"duration,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func (l *Message) String() string {
	line, _ := json.Marshal(l)
	return string(line)
}

// NewCustomLogger creates a new CustomLogger instance and opens a log file.
func NewCustomLogger(logFilePath string) (*CustomLogger, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		return nil, err
	}

	logMulti := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(logMulti, "", log.Ldate|log.Ltime)
	return &CustomLogger{logger, logFile}, nil
}

// Log logs a custom log message with various fields and colored based on log level.
func (c *CustomLogger) Log(lm *Message) {
	logMessage := lm.Level + ": " + lm.String()

	/*
		switch lm.Level {
		case Info:
			logMessage = color.GreenString(logMessage)
		case Warn:
			logMessage = color.YellowString(logMessage)
		case Error:
			logMessage = color.RedString(logMessage)
		}
	*/

	c.Logger.Print(logMessage)
}
