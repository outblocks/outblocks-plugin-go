package plugin

import "github.com/outblocks/outblocks-plugin-go/types"

type PlanRequest struct {
	Apps         []*types.AppPlan        `json:"apps"`
	Dependencies []*types.DependencyPlan `json:"dependencies"`
	TargetApps   []string                `json:"target_apps"`
	Verify       bool                    `json:"verify"`
	Destroy      bool                    `json:"destroy"`

	PluginMap types.PluginStateMap   `json:"plugin_state"`
	Args      map[string]interface{} `json:"args"`
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
