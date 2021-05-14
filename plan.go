package plugin_go

import "github.com/outblocks/outblocks-plugin-go/types"

type PlanState struct {
	Plugin types.PluginState  `json:"plugin"` // plugin name -> object -> state
	Deploy *types.StateDeploy `json:"deploy_state"`
}

type PlanRequest struct {
	Apps         []*types.AppPlanRequest        `json:"apps"`
	Dependencies []*types.DependencyPlanRequest `json:"dependencies"`
	State        *PlanState                     `json:"state"`
	Verify       bool                           `json:"verify"`
}

func (r *PlanRequest) Type() RequestType {
	return RequestTypePlan
}

type PlanResponse struct {
	Deploy *types.Plan `json:"deploy,omitempty"`
	DNS    *types.Plan `json:"dns,omitempty"`
}

func (r *PlanResponse) Type() ResponseType {
	return ResponseTypePlan
}
