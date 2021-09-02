package env

type Enver interface {
	PluginDir() string
	ProjectName() string
	ProjectDir() string
}

var (
	_ Enver = (*Env)(nil)
)
