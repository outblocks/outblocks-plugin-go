package log

import (
	"context"
	"fmt"
	"os"

	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
)

type Log struct {
	cli apiv1.HostServiceClient
}

func NewLogger(cli apiv1.HostServiceClient) Logger {
	return &Log{
		cli: cli,
	}
}

func (l *Log) log(lvl apiv1.LogRequest_Level, msg string) {
	_, _ = l.cli.Log(context.Background(), &apiv1.LogRequest{
		Message: msg,
		Level:   lvl,
	})
}

func (l *Log) writeln(lvl apiv1.LogRequest_Level, a ...interface{}) {
	l.log(lvl, fmt.Sprintln(a...))
}

func (l *Log) writef(lvl apiv1.LogRequest_Level, format string, a ...interface{}) {
	l.log(lvl, fmt.Sprintf(format, a...))
}

func (l *Log) write(lvl apiv1.LogRequest_Level, a ...interface{}) {
	l.log(lvl, fmt.Sprint(a...))
}

func (l *Log) Fatalln(a ...interface{}) {
	l.writeln(apiv1.LogRequest_LEVEL_ERROR, a...)
	os.Exit(1)
}

func (l *Log) Fatalf(format string, a ...interface{}) {
	l.writef(apiv1.LogRequest_LEVEL_ERROR, format, a...)
	os.Exit(1)
}

func (l *Log) Fatal(a ...interface{}) {
	l.write(apiv1.LogRequest_LEVEL_ERROR, a...)
	os.Exit(1)
}

func (l *Log) Errorln(a ...interface{}) {
	l.writeln(apiv1.LogRequest_LEVEL_ERROR, a...)
}

func (l *Log) Errorf(format string, a ...interface{}) {
	l.writef(apiv1.LogRequest_LEVEL_ERROR, format, a...)
}

func (l *Log) Error(a ...interface{}) {
	l.write(apiv1.LogRequest_LEVEL_ERROR, a...)
}

func (l *Log) Warnln(a ...interface{}) {
	l.writeln(apiv1.LogRequest_LEVEL_WARN, a...)
}

func (l *Log) Warnf(format string, a ...interface{}) {
	l.writef(apiv1.LogRequest_LEVEL_WARN, format, a...)
}

func (l *Log) Warn(a ...interface{}) {
	l.write(apiv1.LogRequest_LEVEL_WARN, a...)
}

func (l *Log) Infoln(a ...interface{}) {
	l.writeln(apiv1.LogRequest_LEVEL_INFO, a...)
}

func (l *Log) Infof(format string, a ...interface{}) {
	l.writef(apiv1.LogRequest_LEVEL_INFO, format, a...)
}

func (l *Log) Info(a ...interface{}) {
	l.write(apiv1.LogRequest_LEVEL_INFO, a...)
}

func (l *Log) Debugln(a ...interface{}) {
	l.writeln(apiv1.LogRequest_LEVEL_DEBUG, a...)
}

func (l *Log) Debugf(format string, a ...interface{}) {
	l.writef(apiv1.LogRequest_LEVEL_DEBUG, format, a...)
}

func (l *Log) Debug(a ...interface{}) {
	l.write(apiv1.LogRequest_LEVEL_DEBUG, a...)
}

func (l *Log) Successln(a ...interface{}) {
	l.writeln(apiv1.LogRequest_LEVEL_SUCCESS, a...)
}

func (l *Log) Successf(format string, a ...interface{}) {
	l.writef(apiv1.LogRequest_LEVEL_SUCCESS, format, a...)
}

func (l *Log) Success(a ...interface{}) {
	l.write(apiv1.LogRequest_LEVEL_SUCCESS, a...)
}

func (l *Log) Println(a ...interface{}) {
	l.writeln(apiv1.LogRequest_LEVEL_PRINT, a...)
}

func (l *Log) Printf(format string, a ...interface{}) {
	l.writef(apiv1.LogRequest_LEVEL_PRINT, format, a...)
}

func (l *Log) Print(a ...interface{}) {
	l.write(apiv1.LogRequest_LEVEL_PRINT, a...)
}
