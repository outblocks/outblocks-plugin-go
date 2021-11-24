package env

import "os"

type Env struct{}

func NewEnv() *Env {
	return &Env{}
}

func (e *Env) PluginDir() string {
	return os.Getenv("OUTBLOCKS_PLUGIN_DIR")
}

func (e *Env) PluginProjectCacheDir() string {
	return os.Getenv("OUTBLOCKS_PLUGIN_PROJECT_CACHE_DIR")
}

func (e *Env) ProjectDir() string {
	return os.Getenv("OUTBLOCKS_PROJECT_DIR")
}

func (e *Env) ProjectName() string {
	return os.Getenv("OUTBLOCKS_PROJECT_NAME")
}

func (e *Env) ProjectID() string {
	return os.Getenv("OUTBLOCKS_PROJECT_ID")
}

func (e *Env) Env() string {
	return os.Getenv("OUTBLOCKS_ENV")
}
