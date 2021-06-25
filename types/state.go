package types

import "encoding/json"

type StateSource struct {
	Name    string `json:"name"`
	Created bool   `json:"created"`
}

type DNS struct {
	InternalIP  string `json:"internal_ip"`
	IP          string `json:"ip"`
	InternalURL string `json:"internal_url"`
	URL         string `json:"url"`
}

type AppState struct {
	App *App `json:"app"`
	DNS *DNS `json:"dns"`
}

func NewAppState(app *App) *AppState {
	return &AppState{
		App: app,
	}
}

type DependencyState struct {
	Dependency *Dependency `json:"dependency"`
	DNS        *DNS        `json:"dns"`
}

func NewDependencyState() *DependencyState {
	return &DependencyState{}
}

type PluginStateMap map[string]json.RawMessage

type StateData struct {
	PluginsMap       map[string]PluginStateMap   `json:"plugins_state"` // plugin name -> object -> state
	AppStates        map[string]*AppState        `json:"app_states"`
	DependencyStates map[string]*DependencyState `json:"dep_states"`
}
