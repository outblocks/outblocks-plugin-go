package env

type Enver interface {
	PluginDir() string
	PluginCacheDir() string
	ProjectID() string
	ProjectName() string
	ProjectDir() string
	Env() string
}

var (
	_ Enver = (*Env)(nil)
)
