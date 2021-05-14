package env

type Enver interface {
	PluginDir() string
	ProjectName() string
	ProjectPath() string
}

var (
	_ Enver = (*Env)(nil)
)
