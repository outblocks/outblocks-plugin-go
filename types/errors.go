package types

import (
	"time"

	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	LockErrorMessage       = "lock error"
	StateLockErrorMessage  = "state lock error"
	ValidationErrorMessage = "validation error"
)

func NewStatusValidationError(path, msg string) error {
	st := status.New(codes.InvalidArgument, ValidationErrorMessage)

	st, _ = st.WithDetails(&apiv1.ValidationError{
		Path:    path,
		Message: msg,
	})

	return st.Err()
}

func NewStatusLockError(det ...*apiv1.LockError) error {
	st := status.New(codes.FailedPrecondition, LockErrorMessage)

	for _, d := range det {
		st, _ = st.WithDetails(d)
	}

	return st.Err()
}

func NewLockError(lockname, lockinfo, owner string, createdAt time.Time) *apiv1.LockError {
	return &apiv1.LockError{
		Owner:     owner,
		LockName:  lockname,
		LockInfo:  lockinfo,
		CreatedAt: timestamppb.New(createdAt),
	}
}

func NewStatusStateLockError(lockinfo, owner string, createdAt time.Time) error {
	st := status.New(codes.FailedPrecondition, StateLockErrorMessage)

	st, _ = st.WithDetails(&apiv1.StateLockError{
		Owner:     owner,
		LockInfo:  lockinfo,
		CreatedAt: timestamppb.New(createdAt),
	})

	return st.Err()
}
