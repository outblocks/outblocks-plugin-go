package plugin_go

import "github.com/outblocks/outblocks-plugin-go/types"

type PlanRequest struct {
	Apps         []*types.App        `json:"apps"`
	Dependencies []*types.Dependency `json:"dependencies"`
}

func (r *PlanRequest) Type() RequestType {
	return RequestTypePlan
}

type PlanResponse struct {
	Apps         []*types.AppPlan        `json:"apps,omitempty"`
	Dependencies []*types.DependencyPlan `json:"dependencies,omitempty"`
}

func (r *PlanResponse) Type() ResponseType {
	return ResponseTypePlan
}
