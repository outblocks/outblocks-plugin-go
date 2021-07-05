package log

type Logger interface {
	Fatalln(a ...interface{})
	Fatalf(format string, a ...interface{})
	Errorln(a ...interface{})
	Errorf(format string, a ...interface{})
	Warnln(a ...interface{})
	Warnf(format string, a ...interface{})
	Infoln(a ...interface{})
	Infof(format string, a ...interface{})
	Debugln(a ...interface{})
	Debugf(format string, a ...interface{})
	Successln(a ...interface{})
	Successf(format string, a ...interface{})
}

var (
	_ Logger = (*Log)(nil)
)
