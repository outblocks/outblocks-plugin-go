package plugin

import "github.com/outblocks/outblocks-plugin-go/types"

type DeployBaseRequest struct {
	Apps         []*types.AppPlan        `json:"apps"`
	Dependencies []*types.DependencyPlan `json:"dependencies"`

	Destroy bool `json:"destroy"`

	StateMap types.PluginStateMap   `json:"state_map"`
	Args     map[string]interface{} `json:"args"`
}

type PlanRequest struct {
	DeployBaseRequest
	Verify bool `json:"verify"`
}

func (r *PlanRequest) Type() RequestType {
	return RequestTypePlan
}

type PlanResponse struct {
	DeployPlan *types.Plan `json:"deploy,omitempty"`
	DNSPlan    *types.Plan `json:"dns,omitempty"`

	PluginMap        types.PluginStateMap              `json:"plugin_state"`
	AppStates        map[string]*types.AppState        `json:"app_states"`
	DependencyStates map[string]*types.DependencyState `json:"dep_states"`
}

func (r *PlanResponse) Type() ResponseType {
	return ResponseTypePlan
}
