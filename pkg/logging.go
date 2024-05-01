package pkg

import (
	"log"
	"os"
)

type Logger interface {
	Info(format string, values ...any)
	Error(format string, values ...any)
}

type Log struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

func NewLogger() Log {
	return newLogger(
		log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	)
}

func newLogger(
	infoLogger *log.Logger,
	errorLogger *log.Logger,
) Log {
	return Log{
		InfoLogger:  infoLogger,
		ErrorLogger: errorLogger,
	}
}

func (l Log) Info(format string, values ...any) {
	l.InfoLogger.Printf(format, values...)
}

func (l Log) Error(format string, values ...any) {
	l.ErrorLogger.Printf(format, values...)
}
