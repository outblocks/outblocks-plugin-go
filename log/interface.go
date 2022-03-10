package log

type Logger interface {
	Fatal(a ...interface{})
	Fatalln(a ...interface{})
	Fatalf(format string, a ...interface{})
	Error(a ...interface{})
	Errorln(a ...interface{})
	Errorf(format string, a ...interface{})
	Warn(a ...interface{})
	Warnln(a ...interface{})
	Warnf(format string, a ...interface{})
	Info(a ...interface{})
	Infoln(a ...interface{})
	Infof(format string, a ...interface{})
	Debug(a ...interface{})
	Debugln(a ...interface{})
	Debugf(format string, a ...interface{})
	Success(a ...interface{})
	Successln(a ...interface{})
	Successf(format string, a ...interface{})
	Print(a ...interface{})
	Println(a ...interface{})
	Printf(format string, a ...interface{})
}

var (
	_ Logger = (*Log)(nil)
)
