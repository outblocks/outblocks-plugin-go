package types

import "github.com/outblocks/outblocks-plugin-go/util"

type MonitoringChannelSlack struct {
	Channel string `json:"channel,omitempty"`
	Token   string `json:"token,omitempty"`
}

func NewMonitoringChannelSlack(in map[string]any) (*MonitoringChannelSlack, error) {
	o := &MonitoringChannelSlack{}

	return o, util.MapstructureJSONDecode(in, o)
}

type MonitoringChannelEmail struct {
	Email string `json:"email,omitempty"`
}

func NewMonitoringChannelEmail(in map[string]any) (*MonitoringChannelEmail, error) {
	o := &MonitoringChannelEmail{}

	return o, util.MapstructureJSONDecode(in, o)
}
