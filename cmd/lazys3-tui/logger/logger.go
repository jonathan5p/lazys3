package logger

import (
	"log"
	"os"
)

type Logger struct {
	inner *log.Logger
	file  *os.File
}

func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return &Logger{
		inner: log.New(f, "", log.LstdFlags|log.Lmicroseconds),
		file:  f,
	}, nil
}

func (l *Logger) Close() error {
	return l.file.Close()
}

func (l *Logger) Info(msg string, args ...any) {
	l.inner.Printf("[INFO]  "+msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.inner.Printf("[DEBUG] "+msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.inner.Printf("[ERROR] "+msg, args...)
}
