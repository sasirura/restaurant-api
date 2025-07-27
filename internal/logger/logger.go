// Package logger
package logger

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type Logger struct {
	minLevel log.Level
	mu       sync.Mutex
}

func New(minLevel log.Level, output io.Writer) *Logger {
	log.SetLevel(minLevel)

	if output == nil {
		output = os.Stdout
	}
	log.SetOutput(output)

	return &Logger{
		minLevel: minLevel,
	}
}

func (l *Logger) Info(message string, properties ...interface{}) {
	l.print(log.LevelInfo, message, propertiesToMap(properties))
}

func (l *Logger) Error(message string, properties ...interface{}) {
	l.print(log.LevelError, message, propertiesToMap(properties))
}

func (l *Logger) Fatal(message string, properties ...interface{}) {
	l.print(log.LevelFatal, message, propertiesToMap(properties))
	os.Exit(1)
}

func (l *Logger) Debug(message string, properties ...interface{}) {
	l.print(log.LevelDebug, message, propertiesToMap(properties))
}

func (l *Logger) print(level log.Level, message string, properties map[string]string) {
	if level < l.minLevel {
		return
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Mesaage    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      toString(level),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Mesaage:    message,
		Properties: properties,
	}

	if level >= log.LevelError {
		aux.Trace = string(debug.Stack())
	}

	line, err := json.Marshal(aux)
	if err != nil {
		log.Errorw("Failed to marshal log messages", "error", err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	switch level {
	case log.LevelDebug:
		log.Debug(string(line))
	case log.LevelInfo:
		log.Info(string(line))
	case log.LevelError:
		log.Error(string(line))
	case log.LevelFatal:
		log.Fatal(string(line))
	}
}

func (l *Logger) Write(message []byte) (n int, err error) {
	l.print(log.LevelError, string(message), nil)
	return len(message), nil
}

func propertiesToMap(properties []any) map[string]string {
	result := make(map[string]string)
	for i := 0; i < len(properties)-1; i += 2 {
		key, ok := properties[i].(string)
		if !ok {
			continue
		}
		value, ok := properties[i+1].(string)
		if !ok {
			continue
		}
		result[key] = value
	}
	return result
}

func toString(level log.Level) string {
	switch level {
	case log.LevelDebug:
		return "DEBUG"
	case log.LevelInfo:
		return "INFO"
	case log.LevelError:
		return "ERROR"
	case log.LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}
