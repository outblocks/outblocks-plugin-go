package types

type StateSource struct {
	Name    string `json:"name"`
	Created bool   `json:"created"`
}

type AppState struct {
	IP  string `json:"ip"`
	URL string `json:"url"`
}

type DependencyState struct {
	IP  string `json:"ip"`
	URL string `json:"url"`
}

type PluginState map[string]interface{}

type StateDeploy struct {
	Apps         map[string]*AppState        `json:"apps"`
	Dependencies map[string]*DependencyState `json:"dependencies"`
}

type StateData struct {
	Plugins  map[string]PluginState `json:"plugins"` // plugin name -> object -> state
	Deploy   *StateDeploy           `json:"deploy_state"`
	LockInfo []byte                 `json:"lockinfo"`
}
