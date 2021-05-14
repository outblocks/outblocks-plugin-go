package types

const (
	DNSObject = "dns"
)

type AppPlanRequest struct {
	IsDeploy bool `json:"is_deploy"`
	IsDNS    bool `json:"is_dns"`
	App      *App `json:"app"`
}

type DependencyPlanRequest struct {
	Dependency *Dependency `json:"dependency"`
}

type Plan struct {
	Apps         []*AppPlan        `json:"apps,omitempty"`
	Dependencies []*DependencyPlan `json:"dependencies,omitempty"`
}

type PlanAction struct {
	Object      string `json:"object"`
	Description string `json:"description"`
	Operation   []byte `json:"op"`
}

func (a *PlanAction) IsDNS() bool {
	return a.Object == DNSObject
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
