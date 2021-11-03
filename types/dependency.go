package types

import (
	"fmt"
	"strings"

	"github.com/outblocks/outblocks-plugin-go/util"
)

type Dependency struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

func (d *Dependency) String() string {
	return fmt.Sprintf("Dependency<Name=%s,Type=%s>", d.Name, d.Type)
}

func (d *Dependency) EnvPrefix() string {
	return fmt.Sprintf("DEP_%s_%s_", strings.ToUpper(d.Type), util.SanitizeEnvVar((strings.ToUpper(d.Name))))
}
