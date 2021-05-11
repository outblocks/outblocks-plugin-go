package types

import "fmt"

type Dependency struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Deploy     string                 `json:"deploy"`
	Properties map[string]interface{} `json:"properties"`
}

func (d *Dependency) String() string {
	return fmt.Sprintf("Dependency<Name=%s,Type=%s>", d.Name, d.Type)
}
