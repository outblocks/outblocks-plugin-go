package log

import (
	"fmt"
	"os"
)

type Level byte

const (
	LevelError Level = iota
	LevelWarn
	LevelInfo
	LevelDebug
	LevelSuccess
)

type Log struct{}

func NewLogger() Logger {
	return &Log{}
}

func (l *Log) writeln(lvl Level, a ...interface{}) {
	os.Stderr.Write([]byte{byte(lvl)})
	fmt.Fprintln(os.Stderr, a...)
}

func (l *Log) writef(lvl Level, format string, a ...interface{}) {
	os.Stderr.Write([]byte{byte(lvl)})
	fmt.Fprintf(os.Stderr, format, a...)
}

func (l *Log) Fatalln(a ...interface{}) {
	l.writeln(LevelError, a...)
	os.Exit(1)
}

func (l *Log) Fatalf(format string, a ...interface{}) {
	l.writef(LevelError, format, a...)
	os.Exit(1)
}

func (l *Log) Errorln(a ...interface{}) {
	l.writeln(LevelError, a...)
}

func (l *Log) Errorf(format string, a ...interface{}) {
	l.writef(LevelError, format, a...)
}

func (l *Log) Warnln(a ...interface{}) {
	l.writeln(LevelWarn, a...)
}

func (l *Log) Warnf(format string, a ...interface{}) {
	l.writef(LevelWarn, format, a...)
}

func (l *Log) Infoln(a ...interface{}) {
	l.writeln(LevelInfo, a...)
}

func (l *Log) Infof(format string, a ...interface{}) {
	l.writef(LevelInfo, format, a...)
}

func (l *Log) Debugln(a ...interface{}) {
	l.writeln(LevelDebug, a...)
}

func (l *Log) Debugf(format string, a ...interface{}) {
	l.writef(LevelDebug, format, a...)
}

func (l *Log) Successln(a ...interface{}) {
	l.writeln(LevelSuccess, a...)
}

func (l *Log) Successf(format string, a ...interface{}) {
	l.writef(LevelSuccess, format, a...)
}
