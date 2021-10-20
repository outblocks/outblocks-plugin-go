package fields

import (
	"fmt"

	"github.com/outblocks/outblocks-plugin-go/util"
)

type FieldVarEvaluator struct {
	*util.BaseVarEvaluator
}

func NewFieldVarEvaluator(vars map[string]interface{}) *FieldVarEvaluator {
	return &FieldVarEvaluator{
		BaseVarEvaluator: util.NewBaseVarEvaluator(vars).
			WithEncoder(fieldsVarEncoder).
			WithVarChar('%'),
	}
}

func fieldsVarEncoder(input interface{}) ([]byte, error) {
	switch input.(type) {
	case StringInputField, StringOutputField, string:
		return []byte("%s"), nil
	case IntInputField, IntOutputField, int:
		return []byte("%d"), nil
	}

	return nil, fmt.Errorf("unknown input type")
}

func (e *FieldVarEvaluator) Expand(input string) (StringInputField, error) {
	format, params, err := e.ExpandRaw([]byte(input))
	if err != nil {
		return nil, err
	}

	return Sprintf(string(format), params...), nil
}
