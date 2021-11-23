package types

import (
	"encoding/json"

	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
)

const (
	SourceApp        = "app"
	SourceDependency = "dependency"
	SourcePlugin     = "plugin"
)

func NewPluginState() *apiv1.PluginState {
	return &apiv1.PluginState{
		Other:    make(map[string][]byte),
		Volatile: make(map[string][]byte),
	}
}

func PluginStateFromProto(in *apiv1.PluginState) *PluginState {
	other := make(map[string]json.RawMessage, len(in.Other))
	for k, v := range in.Other {
		other[k] = v
	}

	volatile := make(map[string]json.RawMessage, len(in.Volatile))
	for k, v := range in.Volatile {
		volatile[k] = v
	}

	return &PluginState{
		Registry: in.Registry,
		Other:    other,
		Volatile: volatile,
	}
}

type PluginState struct {
	Registry json.RawMessage            `json:"registry,omitempty"`
	Other    map[string]json.RawMessage `json:"other,omitempty"`
	Volatile map[string]json.RawMessage `json:"volatile,omitempty"`
}

func (p *PluginState) Proto() *apiv1.PluginState {
	other := make(map[string][]byte, len(p.Other))
	for k, v := range p.Other {
		other[k] = v
	}

	volatile := make(map[string][]byte, len(p.Volatile))
	for k, v := range p.Volatile {
		volatile[k] = v
	}

	return &apiv1.PluginState{
		Registry: p.Registry,
		Other:    other,
		Volatile: volatile,
	}
}

type StateData struct {
	Apps         map[string]*apiv1.AppState        `json:"apps"`
	Dependencies map[string]*apiv1.DependencyState `json:"dependencies"`
	Plugins      map[string]*PluginState           `json:"plugins_state"` // plugin name -> object -> state
}
