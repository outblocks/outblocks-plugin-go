package types

import (
	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
)

func NewPluginState() *apiv1.PluginState {
	return &apiv1.PluginState{
		Other: make(map[string][]byte),
	}
}
