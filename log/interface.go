package log

type Logger interface {
	Fatal(a ...any)
	Fatalln(a ...any)
	Fatalf(format string, a ...any)
	Error(a ...any)
	Errorln(a ...any)
	Errorf(format string, a ...any)
	Warn(a ...any)
	Warnln(a ...any)
	Warnf(format string, a ...any)
	Info(a ...any)
	Infoln(a ...any)
	Infof(format string, a ...any)
	Debug(a ...any)
	Debugln(a ...any)
	Debugf(format string, a ...any)
	Success(a ...any)
	Successln(a ...any)
	Successf(format string, a ...any)
	Print(a ...any)
	Println(a ...any)
	Printf(format string, a ...any)
}

var _ Logger = (*Log)(nil)
