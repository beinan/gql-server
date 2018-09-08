package logging

import (
	"log"
	"os"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
}

func StandardLogger(level int) Logger {
	return &stdLogger{level, log.New(os.Stdout, "", log.Ldate|log.Ltime)}
}

type stdLogger struct {
	level  int
	logger *log.Logger
}

func (l *stdLogger) Debug(a ...interface{}) {
	if l.level <= DEBUG {
		l.logger.Println(a...)
	}
}

func (l *stdLogger) Info(a ...interface{}) {
	if l.level <= INFO {
		l.logger.Println(a...)
	}
}

func (l *stdLogger) Error(a ...interface{}) {
	if l.level <= WARN {
		l.logger.Println(a...)
	}
}

func (l *stdLogger) Warn(a ...interface{}) {
	if l.level <= ERROR {
		l.logger.Println(a...)
	}
}
