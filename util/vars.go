package util

import (
	"bytes"
	"fmt"
	"strings"
)

type BaseVarEvaluator struct {
	vars                  map[string]interface{}
	encoder               func(val interface{}) ([]byte, error)
	keyGetter             func(vars map[string]interface{}, key string) (val interface{}, ok bool)
	ignoreComments        bool
	varChar, commentsChar byte
}

func NewBaseVarEvaluator(vars map[string]interface{}) *BaseVarEvaluator {
	return &BaseVarEvaluator{
		vars:           vars,
		keyGetter:      defaultVarKeyGetter,
		encoder:        defaultVarEncoder,
		ignoreComments: false,
		varChar:        '$',
		commentsChar:   '#',
	}
}

func (e *BaseVarEvaluator) WithEncoder(encoder func(val interface{}) ([]byte, error)) *BaseVarEvaluator {
	e.encoder = encoder
	return e
}

func (e *BaseVarEvaluator) WithKeyGetter(keyGetter func(vars map[string]interface{}, key string) (val interface{}, ok bool)) *BaseVarEvaluator {
	e.keyGetter = keyGetter
	return e
}

func (e *BaseVarEvaluator) WithIgnoreComments(ignoreComments bool) *BaseVarEvaluator {
	e.ignoreComments = ignoreComments
	return e
}

func (e *BaseVarEvaluator) WithVarChar(varChar byte) *BaseVarEvaluator {
	e.varChar = varChar
	return e
}

func (e *BaseVarEvaluator) WithCommentsChar(commentsChar byte) *BaseVarEvaluator {
	e.commentsChar = commentsChar
	return e
}

func defaultVarEncoder(input interface{}) ([]byte, error) {
	return []byte("%v"), nil
}

func defaultVarKeyGetter(vars map[string]interface{}, key string) (val interface{}, ok bool) {
	parts := strings.Split(key, ".")

	for _, part := range parts[:len(parts)-1] {
		vars, ok = vars[part].(map[string]interface{})
		if !ok {
			return nil, false
		}
	}

	val, ok = vars[parts[len(parts)-1]]

	return val, ok
}

func (e *BaseVarEvaluator) ExpandRaw(input []byte) (output []byte, params []interface{}, err error) {
	tokenStart := -1

	var token string

	in := bytes.Split(bytes.ReplaceAll(input, []byte{'\r', '\n'}, []byte{'\n'}), []byte{'\n'})
	out := make([][]byte, len(in))

	for l, line := range in {
		ll := len(line)
		done := 0

		if e.ignoreComments {
			lineTrimmed := bytes.TrimSpace(line)
			if len(lineTrimmed) > 0 && lineTrimmed[0] == e.commentsChar {
				out[l] = line

				continue
			}
		}

		for c := range line {
			switch {
			case tokenStart == -1:
				// todo: magic marker set
				if c+1 < ll && line[c] == e.varChar && line[c+1] == '{' {
					c++
					tokenStart = c
				}

			case line[c] == '}':
				token = string(line[tokenStart+1 : c])
				if token == "" {
					return nil, nil, fmt.Errorf("[%d:%d] empty expansion found", l+1, tokenStart)
				}

				out[l] = append(out[l], line[done:tokenStart-1]...)

				val, ok := e.keyGetter(e.vars, token)
				if !ok {
					return nil, nil, fmt.Errorf("[%d:%d] expansion value '%s' could not be evaluated", l+1, tokenStart, token)
				}

				valOut, err := e.encoder(val)
				if err != nil {
					return nil, nil, fmt.Errorf("[%d:%d] expansion value '%s' could not be encoded, unknown field\nerror: %w",
						l+1, tokenStart, token, err)
				}

				out[l] = append(out[l], valOut...)
				params = append(params, val)

				done = c + 1
				tokenStart = -1
			}
		}

		out[l] = append(out[l], line[done:]...)
	}

	return bytes.Join(out, []byte{'\n'}), params, nil
}
