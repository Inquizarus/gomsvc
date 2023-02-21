package logging

import (
	"io"
	"log"
)

type plainLogger struct {
	log *log.Logger
}

func (l *plainLogger) print(level string, args ...interface{}) {
	for i := 0; i < len(args); i++ {
		entry := args[i]
		l.log.Printf("[%s] %s", level, entry)
	}
}

func (l *plainLogger) Info(args ...interface{}) {
	l.print("INFO", args)
}
func (l *plainLogger) Infof(format string, args ...interface{}) {}
func (l *plainLogger) Debug(args ...interface{}) {
	l.print("DEBUG", args)
}
func (l *plainLogger) Debugf(format string, args ...interface{}) {}
func (l *plainLogger) Error(args ...interface{})                 {}
func (l *plainLogger) Errorf(format string, args ...interface{}) {}

func NewPlainLogger(out io.Writer, prefix string) Logger {
	log := log.New(out, prefix, 0)
	return &plainLogger{
		log,
	}
}
