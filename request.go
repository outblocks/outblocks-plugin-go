package plugin

type RequestType int

const (
	RequestTypeInit RequestType = iota + 1
	RequestTypeStart
	RequestTypePlan
	RequestTypeApply
	RequestTypeRun
	RequestTypeGetState
	RequestTypeSaveState
	RequestTypeReleaseLock
	RequestTypeCommand
	RequestTypePromptAnswer
)

type Request interface {
	Type() RequestType
}

var (
	_ Request = (*InitRequest)(nil)
	_ Request = (*StartRequest)(nil)
	_ Request = (*ApplyRequest)(nil)
	_ Request = (*PlanRequest)(nil)
	_ Request = (*GetStateRequest)(nil)
	_ Request = (*SaveStateRequest)(nil)
)

type RequestHeader struct {
	Type RequestType `json:"type"`
}
