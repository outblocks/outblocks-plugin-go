package plugin_go

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
type ForceUnlockRequest struct {
	LockInfo   string                 `json:"lockinfo"`
	StateType  string                 `json:"type"`
	Env        string                 `json:"env"`
	Properties map[string]interface{} `json:"properties"`
}

func (r *ForceUnlockRequest) Type() RequestType {
	return RequestTypeForceUnlock
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
	LockInfo   string                 `json:"lockinfo"`
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
