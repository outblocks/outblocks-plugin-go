package plugin

import "github.com/outblocks/outblocks-plugin-go/types"

type ApplyRequest struct {
	Apps         []*types.AppPlan        `json:"apps"`
	Dependencies []*types.DependencyPlan `json:"dependencies"`
	TargetApps   []string                `json:"target_apps"`
	SkipApps     []string                `json:"skip_apps"`

	Destroy bool `json:"destroy"`

	PluginMap types.PluginStateMap   `json:"plugin_state"`
	Args      map[string]interface{} `json:"args"`
}

func (r *ApplyRequest) Type() RequestType {
	return RequestTypeApply
}

type ApplyResponse struct {
	Actions []*types.ApplyAction `json:"actions"`
}

func (r *ApplyResponse) Type() ResponseType {
	return ResponseTypeApply
}

type ApplyDoneResponse struct {
	PluginMap        types.PluginStateMap              `json:"plugin_state"`
	AppStates        map[string]*types.AppState        `json:"app_states"`
	DependencyStates map[string]*types.DependencyState `json:"dep_states"`
}

func (r *ApplyDoneResponse) Type() ResponseType {
	return ResponseTypeApplyDone
}
