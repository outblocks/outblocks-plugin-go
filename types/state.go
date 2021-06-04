package types

type StateSource struct {
	Name    string `json:"name"`
	Created bool   `json:"created"`
}

type DeployDNS struct {
	InternalIP  string `json:"internal_ip"`
	ExternalIP  string `json:"external_ip"`
	InternalURL string `json:"internal_url"`
	ExternalURL string `json:"external_url"`
}

type DNS struct {
	IP  string `json:"ip"`
	URL string `json:"url"`
}

type AppState struct {
	App         *App                   `json:"app"`
	DeployState map[string]interface{} `json:"deploy_state"`
	DeployDNS   *DeployDNS             `json:"deploy_dns"`
	DNSState    map[string]interface{} `json:"dns_state"`
	DNS         *DNS                   `json:"dns"`
}

func NewAppState() *AppState {
	return &AppState{
		DeployState: make(map[string]interface{}),
		DNSState:    make(map[string]interface{}),
	}
}

type DependencyState struct {
	Dependency  *Dependency            `json:"dependency"`
	DeployState map[string]interface{} `json:"deploy_state"`
	DeployDNS   *DeployDNS             `json:"deploy_dns"`
	DNSState    map[string]interface{} `json:"dns_state"`
	DNS         *DNS                   `json:"dns"`
}

func NewDependencyState() *DependencyState {
	return &DependencyState{
		DeployState: make(map[string]interface{}),
		DNSState:    make(map[string]interface{}),
	}
}

type PluginStateMap map[string]interface{}

type StateData struct {
	PluginsMap       map[string]PluginStateMap   `json:"plugins_state"` // plugin name -> object -> state
	AppStates        map[string]*AppState        `json:"app_states"`
	DependencyStates map[string]*DependencyState `json:"dep_states"`
}
