package communication

type ResponseType int

const (
	ResponseTypeEmpty ResponseType = iota + 1
	ResponseTypeData
	ResponseTypeMap
	ResponseTypePrompt
	ResponseTypeMessage
	ResponseTypeError
	ResponseTypeValidationError
	ResponseTypeUnhandled
)

type Response interface {
	Type() ResponseType
}

type StreamingResponse interface {
	Response
	IsFinal() bool
}

var (
	_ Response = (*GetStateResponse)(nil)
	_ Response = (*SaveStateResponse)(nil)
	_ Response = (*UnhandledResponse)(nil)
)

type ResponseHeader struct {
	Type        ResponseType `json:"type"`
	IsStreaming bool         `json:"is_streaming"`
	IsFinal     bool         `json:"is_final"`
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

type MessageResponse struct {
	LogLevel string `json:"level"`
	Message  string `json:"message"`
}

func (r *MessageResponse) Level() string {
	if r.LogLevel != "" {
		return r.LogLevel
	}

	return "info"
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

type MapResponse struct {
	Map map[string]interface{} `json:"data"`
}

func (r *MapResponse) Type() ResponseType {
	return ResponseTypeMap
}

func (r *MapResponse) Data() map[string]interface{} {
	return r.Map
}

type EmptyResponse struct{}

func (r *EmptyResponse) Type() ResponseType {
	return ResponseTypeEmpty
}

type PromptResponse struct {
	// TODO: prompt support
}

func (r *PromptResponse) Type() ResponseType {
	return ResponseTypePrompt
}

type ValidationErrorResponse struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

func (r *ValidationErrorResponse) Type() ResponseType {
	return ResponseTypeValidationError
}
