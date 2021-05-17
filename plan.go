package plugin_go

import "github.com/outblocks/outblocks-plugin-go/types"

type PlanRequest struct {
	Apps         []*types.AppInfo        `json:"apps"`
	Dependencies []*types.DependencyInfo `json:"dependencies"`
	PluginState  types.PluginStateMap    `json:"plugin_state"`
	Verify       bool                    `json:"verify"`
}

func (r *PlanRequest) Type() RequestType {
	return RequestTypePlan
}

type PlanResponse struct {
	DeployPlan *types.Plan `json:"deploy,omitempty"`
	DNSPlan    *types.Plan `json:"dns,omitempty"`
}

func (r *PlanResponse) Type() ResponseType {
	return ResponseTypePlan
}
