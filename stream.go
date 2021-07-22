package plugin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

type ReceiverStream struct {
	conn io.ReadWriteCloser
	r    *bufio.Reader
}

func NewReceiverStream(c io.ReadWriteCloser) *ReceiverStream {
	return &ReceiverStream{
		conn: c,
		r:    bufio.NewReader(c),
	}
}

func (s *ReceiverStream) Send(res Response) error {
	return writeResponse(s.conn, res)
}

func (s *ReceiverStream) Recv() (Request, error) {
	return readRequest(s.r)
}

func (s *ReceiverStream) Close() error {
	return s.conn.Close()
}

func writeResponse(w io.Writer, res Response) error {
	header := &ResponseHeader{
		Type: res.Type(),
	}

	// Prepare header.
	headerData, err := json.Marshal(header)
	if err != nil {
		return err
	}

	// Prepare response itself.
	responseData, err := json.Marshal(res)
	if err != nil {
		return err
	}

	data := headerData
	data = append(data, byte('\n'))
	data = append(data, responseData...)
	data = append(data, byte('\n'))

	_, err = w.Write(data)

	return err
}

func readRequest(r *bufio.Reader) (Request, error) {
	data, err := r.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("unable to read header: %w", err)
	}

	var header RequestHeader
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, fmt.Errorf("header unmarshal error: %s", err)
	}

	var req Request

	switch header.Type {
	case RequestTypeInit:
		req = &InitRequest{}
	case RequestTypeStart:
		req = &StartRequest{}
	case RequestTypePlan:
		req = &PlanRequest{}
	case RequestTypeApply:
		req = &ApplyRequest{}
	case RequestTypeRun:
		req = &RunRequest{}
	case RequestTypeGetState:
		req = &GetStateRequest{}
	case RequestTypeSaveState:
		req = &SaveStateRequest{}
	case RequestTypeCommand:
		req = &CommandRequest{}
	case RequestTypeReleaseLock:
		req = &ReleaseLockRequest{}
	case RequestTypePromptAnswer:
		req = &PromptInputAnswer{}
	case RequestTypePromptConfirmation:
		req = &PromptConfirmationAnswer{}
	default:
		panic(fmt.Sprintf("unknown request type: %d", header.Type))
	}

	data, err = r.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("unable to read request: %w", err)
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("request unmarshal error: %s", err)
	}

	return req, nil
}
