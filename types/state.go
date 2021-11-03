package types

import "encoding/json"

type StateSource struct {
	Name    string `json:"name"`
	Created bool   `json:"created"`
}

type SSLStatus string

const (
	SSLStatusUnknown            SSLStatus = "UNKNOWN"
	SSLStatusOK                 SSLStatus = "OK"
	SSLStatusProvisioning       SSLStatus = "PROVISIONING"
	SSLStatusProvisioningFailed SSLStatus = "PROVISIONING FAILED"
	SSLStatusRenewalFailed      SSLStatus = "RENEWAL FAILED"
)

type DNSState struct {
	InternalIP     string                 `json:"internal_ip,omitempty"`
	IP             string                 `json:"ip"`
	CNAME          string                 `json:"cname,omitempty"`
	InternalURL    string                 `json:"internal_url,omitempty"`
	URL            string                 `json:"url"`
	Manual         bool                   `json:"manual"`
	SSLStatus      SSLStatus              `json:"ssl_status"`
	SSLStatusInfo  string                 `json:"ssl_status_info"`
	ConnectionInfo string                 `json:"connection_info"`
	Properties     map[string]interface{} `json:"properties"`
}

type DeploymentState struct {
	Ready   bool   `json:"ready"`
	Message string `json:"message"`
}

type AppState struct {
	App *App `json:"app"`

	Deployment *DeploymentState `json:"deploy_state"`
	DNS        *DNSState        `json:"dns_state"`
}

func NewAppState(app *App) *AppState {
	return &AppState{
		App: app,
	}
}

type DependencyState struct {
	Dependency *Dependency `json:"dependency"`
	DNS        *DNSState   `json:"dns"`
}

func NewDependencyState(dep *Dependency) *DependencyState {
	return &DependencyState{
		Dependency: dep,
	}
}

type PluginStateMap map[string]json.RawMessage

type StateData struct {
	Apps         map[string]*App           `json:"apps"`
	Dependencies map[string]*Dependency    `json:"dependencies"`
	PluginsMap   map[string]PluginStateMap `json:"plugins_state"` // plugin name -> object -> state
}
