package env

type Enver interface {
	PluginDir() string
	PluginProjectCacheDir() string
	ProjectID() string
	ProjectName() string
	ProjectDir() string
	Env() string
}

var _ Enver = (*Env)(nil)
