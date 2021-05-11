package types

const (
	DNSObject = "dns"
)

type PlanAction struct {
	Object      string `json:"object"`
	Description string `json:"description"`
	Operation   []byte `json:"op"`
}

type AppPlan struct {
	App    *App          `json:"app"`
	Add    []*PlanAction `json:"add"`
	Change []*PlanAction `json:"change"`
	Remove []*PlanAction `json:"remove"`
}

type DependencyPlan struct {
	Dependency *Dependency   `json:"dependency"`
	Add        []*PlanAction `json:"add"`
	Change     []*PlanAction `json:"change"`
	Remove     []*PlanAction `json:"remove"`
}
