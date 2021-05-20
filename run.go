package plugin

import "github.com/outblocks/outblocks-plugin-go/types"

type RunRequest struct {
	Apps         []*types.App        `json:"apps"`
	Dependencies []*types.Dependency `json:"dependencies"`
}

func (r *RunRequest) Type() RequestType {
	return RequestTypeRun
}
