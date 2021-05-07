package log

import (
	"fmt"
	"os"
)

type Log struct{}

func NewLogger() *Log {
	return &Log{}
}

func (l *Log) Fatalln(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func (l *Log) Fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func (l *Log) Errorln(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func (l *Log) Errorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}
