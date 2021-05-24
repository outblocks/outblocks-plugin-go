package plugin

import "github.com/outblocks/outblocks-plugin-go/types"

type PlanRequest struct {
	Apps         []*types.AppPlan        `json:"apps"`
	Dependencies []*types.DependencyPlan `json:"dependencies"`
	PluginState  types.PluginStateMap    `json:"plugin_state"`
	Verify       bool                    `json:"verify"`
	Destroy      bool                    `json:"destroy"`
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
