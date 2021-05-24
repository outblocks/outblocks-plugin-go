package plugin

import (
	"fmt"

	"github.com/outblocks/outblocks-plugin-go/types"
)

// Get State.
type GetStateRequest struct {
	StateType  string                 `json:"type"`
	Env        string                 `json:"env"`
	Properties map[string]interface{} `json:"properties"`
	Lock       bool                   `json:"lock"`
}

func (r *GetStateRequest) Type() RequestType {
	return RequestTypeGetState
}

type GetStateResponse struct {
	State    *types.StateData   `json:"state"`
	LockInfo string             `json:"lockinfo"`
	Source   *types.StateSource `json:"source"`
}

func (r *GetStateResponse) Type() ResponseType {
	return ResponseTypeGetState
}

// Force Unlock.
type ReleaseLockRequest struct {
	LockID     string                 `json:"lock_id"`
	StateType  string                 `json:"type"`
	Env        string                 `json:"env"`
	Properties map[string]interface{} `json:"properties"`
}

func (r *ReleaseLockRequest) Type() RequestType {
	return RequestTypeReleaseLock
}

// Lock Error.
type LockErrorResponse struct {
	Owner    string `json:"owner"`
	LockInfo string `json:"lockinfo"`
}

func (r *LockErrorResponse) Type() ResponseType {
	return ResponseTypeLockError
}

func (r *LockErrorResponse) Error() string {
	return fmt.Sprintf("state lock already acquired by %s", r.Owner)
}

// Save State.
type SaveStateRequest struct {
	State      *types.StateData       `json:"state"`
	StateType  string                 `json:"type"`
	Env        string                 `json:"env"`
	Properties map[string]interface{} `json:"properties"`
}

func (r *SaveStateRequest) Type() RequestType {
	return RequestTypeSaveState
}

type SaveStateResponse struct {
	EmptyResponse
}

func (r *SaveStateResponse) Type() ResponseType {
	return ResponseTypeSaveState
}
