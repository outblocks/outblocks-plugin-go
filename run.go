package plugin

import "github.com/outblocks/outblocks-plugin-go/types"

type RunRequest struct {
	Apps         []*types.AppRun        `json:"apps,omitempty"`
	Dependencies []*types.DependencyRun `json:"dependencies,omitempty"`
	Hosts        map[string]string      `json:"hosts,omitempty"`
	Args         map[string]interface{} `json:"args"`
}

func (r *RunRequest) Type() RequestType {
	return RequestTypeRun
}

type RunningResponse struct {
	Vars map[string]map[string]string // ID -> key->val
}

func (r *RunningResponse) Type() ResponseType {
	return ResponseTypeRunning
}

type RunOutpoutSource int

const (
	RunOutpoutSourceApp RunOutpoutSource = iota + 1
	RunOutpoutSourceDependency
)

type RunOutputResponse struct {
	Source   RunOutpoutSource `json:"source"`
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	IsStderr bool             `json:"is_stderr"`
	Message  string           `json:"message"`
}

func (r *RunOutputResponse) Type() ResponseType {
	return ResponseTypeRunOutput
}
