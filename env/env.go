package env

import "os"

type Env struct{}

func NewEnv() *Env {
	return &Env{}
}

func (e *Env) PluginDir() string {
	return os.Getenv("OUTBLOCKS_PLUGIN_DIR")
}

func (e *Env) ProjectPath() string {
	return os.Getenv("OUTBLOCKS_PROJECT_PATH")
}
