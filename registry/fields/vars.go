package fields

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/outblocks/outblocks-plugin-go/util"
)

var escapePercent = regexp.MustCompile(`%([^{]|$)`)

type FieldVarEvaluator struct {
	*util.BaseVarEvaluator
}

func NewFieldVarEvaluator(vars map[string]any) *FieldVarEvaluator {
	return &FieldVarEvaluator{
		BaseVarEvaluator: util.NewBaseVarEvaluator(vars).
			WithEncoder(fieldsVarEncoder).
			WithVarChar('%').
			WithIgnoreInvalid(true).
			WithSkipRowColumnInfo(true),
	}
}

func fieldsVarEncoder(c *util.VarContext, input any) ([]byte, error) {
	switch input.(type) {
	case StringInputField, StringOutputField, string:
		return []byte("%s"), nil
	case IntInputField, IntOutputField, int:
		return []byte("%d"), nil
	}

	return nil, fmt.Errorf("unknown input type")
}

func (e *FieldVarEvaluator) Expand(input string) (StringInputField, error) {
	input = escapePercent.ReplaceAllString(input, "%$0")

	format, params, err := e.ExpandRaw([]byte(input))
	if err != nil {
		return nil, err
	}

	return Sprintf(strings.ReplaceAll(string(format), "%{", "%%{"), params...), nil
}
