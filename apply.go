package communication

import "github.com/outblocks/outblocks-plugin-go/types"

type ApplyRequest struct {
	Apps         []*types.App           `json:"apps"`
	Dependencies []*types.Dependency    `json:"dependencies"`
	Plan         map[string]interface{} `json:"plan"`
}

func (r *ApplyRequest) Type() RequestType {
	return RequestTypeApply
}
