package env

type Enver interface {
	PluginDir() string
	ProjectPath() string
}

var (
	_ Enver = (*Env)(nil)
)
