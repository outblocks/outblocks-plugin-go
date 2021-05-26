package types

import "fmt"

type App struct {
	Name       string                 `json:"name"`
	Path       string                 `json:"path"`
	Type       string                 `json:"type"`
	URL        string                 `json:"url"`
	Deploy     string                 `json:"deploy"`
	Needs      map[string]*AppNeed    `json:"needs"`
	Properties map[string]interface{} `json:"properties"`
}

func (a *App) TargetName() string {
	return fmt.Sprintf("App '%s'", a.Name)
}

func (a *App) String() string {
	return fmt.Sprintf("App<Name=%s,Type=%s>", a.Name, a.Type)
}

type AppNeed struct {
	Dependency string                 `json:"dependency"`
	Properties map[string]interface{} `json:"properties"`
}
