package types

import (
	"time"

	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewValidationError(path, msg string) error {
	st := status.New(codes.InvalidArgument, "validation error")

	st, _ = st.WithDetails(&apiv1.ValidationError{
		Path:    path,
		Message: msg,
	})

	return st.Err()
}

func NewLockError(lockname, lockinfo, owner string, createdAt time.Time) error {
	st := status.New(codes.FailedPrecondition, "lock error")

	st, _ = st.WithDetails(&apiv1.LockError{
		Owner:     owner,
		LockName:  lockname,
		LockInfo:  lockinfo,
		CreatedAt: timestamppb.New(createdAt),
	})

	return st.Err()
}

func NewStateLockError(lockinfo, owner string, createdAt time.Time) error {
	st := status.New(codes.FailedPrecondition, "state lock error")

	st, _ = st.WithDetails(&apiv1.StateLockError{
		Owner:     owner,
		LockInfo:  lockinfo,
		CreatedAt: timestamppb.New(createdAt),
	})

	return st.Err()
}
