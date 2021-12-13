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

type DNSRecordKey struct {
	Record string
	Type   apiv1.DNSRecord_Type
}

type DNSRecordValue struct {
	Value   string
	Created bool
}

type DNSRecordMap map[DNSRecordKey]DNSRecordValue

func (m DNSRecordMap) MarshalJSON() ([]byte, error) {
	out := make([]*apiv1.DNSRecord, 0, len(m))

	for k, v := range m {
		out = append(out, &apiv1.DNSRecord{
			Record:  k.Record,
			Type:    k.Type,
			Value:   v.Value,
			Created: v.Created,
		})
	}

	return json.Marshal(out)
}

func (m DNSRecordMap) UnmarshalJSON(b []byte) error {
	var out []*apiv1.DNSRecord

	err := json.Unmarshal(b, &out)
	if err != nil {
		return err
	}

	m = make(DNSRecordMap)

	for _, v := range out {
		m[DNSRecordKey{
			Record: v.Record,
			Type:   v.Type,
		}] = DNSRecordValue{
			Value:   v.Value,
			Created: v.Created,
		}
	}

	return nil
}

type StateData struct {
	Apps         map[string]*apiv1.AppState        `json:"apps"`
	Dependencies map[string]*apiv1.DependencyState `json:"dependencies"`
	Plugins      map[string]*PluginState           `json:"plugins_state"` // plugin name -> object -> state

	DNSRecords DNSRecordMap `json:"dns_records"`
}

func (d *StateData) Reset() {
	d.Apps = make(map[string]*apiv1.AppState)
	d.Dependencies = make(map[string]*apiv1.DependencyState)
	d.DNSRecords = make(DNSRecordMap)
}

func (d *StateData) AddDNSRecord(v *apiv1.DNSRecord) {
	d.DNSRecords[DNSRecordKey{
		Record: v.Record,
		Type:   v.Type,
	}] = DNSRecordValue{
		Value:   v.Value,
		Created: v.Created,
	}
}
