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
	RequestTypeReleaseStateLock
	RequestTypeCommand
	RequestTypePromptConfirmation
	RequestTypePromptAnswer
	RequestTypeAcquireLocks
	RequestTypeReleaseLocks
)

type Request interface {
	Type() RequestType
}

var (
	_ Request = (*ProjectInitRequest)(nil)
	_ Request = (*StartRequest)(nil)
	_ Request = (*ApplyRequest)(nil)
	_ Request = (*PlanRequest)(nil)
	_ Request = (*GetStateRequest)(nil)
	_ Request = (*SaveStateRequest)(nil)
)

type RequestHeader struct {
	Type RequestType `json:"type"`
}
