package plugin

type ResponseType int

const (
	ResponseTypeEmpty ResponseType = iota + 1
	ResponseTypeData
	ResponseTypeGetState
	ResponseTypeSaveState
	ResponseTypeLockError
	ResponseTypePlan
	ResponseTypeApply
	ResponseTypeApplyDone
	ResponseTypeMessage
	ResponseTypeError
	ResponseTypeValidationError
	ResponseTypeUnhandled
	ResponseTypeInit
	ResponseTypePromptConfirmation
	ResponseTypePromptSelect
	ResponseTypePromptInput
	ResponseTypeRunning
	ResponseTypeRunOutput
)

type Response interface {
	Type() ResponseType
}

var (
	_ Response = (*GetStateResponse)(nil)
	_ Response = (*SaveStateResponse)(nil)
	_ Response = (*UnhandledResponse)(nil)
)

type ResponseHeader struct {
	Type ResponseType `json:"type"`
}

type UnhandledResponse struct {
}

func (r *UnhandledResponse) Type() ResponseType {
	return ResponseTypeUnhandled
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (r *ErrorResponse) Type() ResponseType {
	return ResponseTypeError
}

type MessageLogLevel string

const (
	MessageLogLevelError   MessageLogLevel = "error"
	MessageLogLevelWarn    MessageLogLevel = "warn"
	MessageLogLevelInfo    MessageLogLevel = "info"
	MessageLogLevelDebug   MessageLogLevel = "debug"
	MessageLogLevelSuccess MessageLogLevel = "success"
)

type MessageResponse struct {
	LogLevel MessageLogLevel `json:"level"`
	Message  string          `json:"message"`
}

func (r *MessageResponse) Level() MessageLogLevel {
	if r.LogLevel != "" {
		return r.LogLevel
	}

	return MessageLogLevelError
}

func (r *MessageResponse) Type() ResponseType {
	return ResponseTypeMessage
}

type DataResponse struct {
	Content []byte `json:"data"`
}

func (r *DataResponse) Type() ResponseType {
	return ResponseTypeData
}

func (r *DataResponse) Data() []byte {
	return r.Content
}

type EmptyResponse struct{}

func (r *EmptyResponse) Type() ResponseType {
	return ResponseTypeEmpty
}

type ValidationErrorResponse struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

func (r *ValidationErrorResponse) Type() ResponseType {
	return ResponseTypeValidationError
}
