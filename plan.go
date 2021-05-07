package communication

import "github.com/outblocks/outblocks-plugin-go/types"

type PlanRequest struct {
	Apps         []*types.App        `json:"apps"`
	Dependencies []*types.Dependency `json:"dependencies"`
}

func (r *PlanRequest) Type() RequestType {
	return RequestTypePlan
}
