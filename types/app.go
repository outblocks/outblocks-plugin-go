package types

import (
	"fmt"
	"strings"

	"github.com/outblocks/outblocks-plugin-go/util"
)

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

func (a *App) EnvPrefix() string {
	return fmt.Sprintf("APP_%s_%s_", strings.ToUpper(a.Type), util.SanitizeEnvVar(strings.ToUpper(a.Name)))
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
