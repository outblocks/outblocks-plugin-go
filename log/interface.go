package log

type Logger interface {
	Fatalln(a ...interface{})
	Fatalf(format string, a ...interface{})
	Errorln(a ...interface{})
	Errorf(format string, a ...interface{})
}

var (
	_ Logger = (*Log)(nil)
)
