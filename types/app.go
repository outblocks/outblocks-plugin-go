package types

import "fmt"

type App struct {
	ID           string                 `json:"id"`
	Dir          string                 `json:"dir"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	URL          string                 `json:"url"`
	PathRedirect string                 `json:"path_redirect"`
	Needs        map[string]*AppNeed    `json:"needs"`
	Properties   map[string]interface{} `json:"properties"`
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
