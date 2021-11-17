package plugin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/outblocks/outblocks-plugin-go/types"
)

// Get State.
type GetStateRequest struct {
	StateType  string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Lock       bool                   `json:"lock"`
	LockWait   time.Duration          `json:"lock_wait"`
}

func (r *GetStateRequest) Type() RequestType {
	return RequestTypeGetState
}

type GetStateResponse struct {
	State    json.RawMessage    `json:"state"`
	LockInfo string             `json:"lock_info"`
	Source   *types.StateSource `json:"source"`
}

func (r *GetStateResponse) Type() ResponseType {
	return ResponseTypeGetState
}

// Release state lock.
type ReleaseStateLockRequest struct {
	LockInfo   string                 `json:"lock_info"`
	StateType  string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

func (r *ReleaseStateLockRequest) Type() RequestType {
	return RequestTypeReleaseStateLock
}

// Locking.
type AcquireLocksRequest struct {
	LockNames  []string               `json:"lock_names,omitempty"`
	LockWait   time.Duration          `json:"lock_wait"`
	Properties map[string]interface{} `json:"properties"`
}

func (r *AcquireLocksRequest) Type() RequestType {
	return RequestTypeAcquireLocks
}

type ReleaseLocksRequest struct {
	Locks      map[string]string      `json:"locks"`
	Properties map[string]interface{} `json:"properties"`
}

func (r *ReleaseLocksRequest) Type() RequestType {
	return RequestTypeReleaseLocks
}

// Lock Error.
type LockErrorResponse struct {
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
	LockInfo  string    `json:"lock_info"`
}

func (r *LockErrorResponse) Type() ResponseType {
	return ResponseTypeLockError
}

func (r *LockErrorResponse) Error() string {
	return fmt.Sprintf("state lock already acquired by %s", r.Owner)
}

// Locks Acquired.
type LocksAcquiredResponse struct {
	LockInfo []string `json:"lock_info"`
}

func (r *LocksAcquiredResponse) Type() ResponseType {
	return ResponseTypeLocksAcquired
}

// Save State.
type SaveStateRequest struct {
	State      json.RawMessage        `json:"state"`
	StateType  string                 `json:"type"`
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
