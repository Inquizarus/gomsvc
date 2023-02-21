package logging

import "github.com/sirupsen/logrus"

type logrusLogger struct {
	log    *logrus.Logger
	fields map[string]interface{}
}

type Fields map[string]interface{}

func (l *logrusLogger) Info(args ...interface{}) {
	l.log.WithFields(l.fields).Info(args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.log.WithFields(l.fields).Infof(format, args...)
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.log.WithFields(l.fields).Debug(args...)
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.log.WithFields(l.fields).Debugf(format, args...)
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.log.WithFields(l.fields).Error(args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.log.WithFields(l.fields).Errorf(format, args...)
}

func NewLogrusLogger(log *logrus.Logger, level string, fields Fields) Logger {
	if nil == log {
		log = logrus.New()
		log.SetFormatter(&logrus.JSONFormatter{})
		logrusLevel, _ := logrus.ParseLevel(level)
		log.SetLevel(logrusLevel)
	}
	return &logrusLogger{
		log,
		fields,
	}
}
