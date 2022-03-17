package util

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var validKey = regexp.MustCompile(`^[a-z]+(\.[a-zA-Z0-9_-]+)*$`)

type VarContext struct {
	Input            []byte
	Line             []byte
	Row              int
	Token            string
	TokenColumnStart int
	TokenColumnEnd   int
}

type BaseVarEvaluator struct {
	vars                  map[string]interface{}
	encoder               func(c *VarContext, val interface{}) ([]byte, error)
	keyGetter             func(c *VarContext, vars map[string]interface{}) (val interface{}, err error)
	ignoreComments        bool
	ignoreInvalid         bool
	skipRowColumnInfo     bool
	varChar, commentsChar byte
}

func NewBaseVarEvaluator(vars map[string]interface{}) *BaseVarEvaluator {
	return &BaseVarEvaluator{
		vars:              vars,
		keyGetter:         DefaultVarKeyGetter,
		encoder:           DefaultVarEncoder,
		ignoreComments:    false,
		ignoreInvalid:     false,
		skipRowColumnInfo: false,
		varChar:           '$',
		commentsChar:      '#',
	}
}

func (e *BaseVarEvaluator) WithEncoder(encoder func(c *VarContext, val interface{}) ([]byte, error)) *BaseVarEvaluator {
	e.encoder = encoder
	return e
}

func (e *BaseVarEvaluator) WithKeyGetter(keyGetter func(c *VarContext, vars map[string]interface{}) (val interface{}, err error)) *BaseVarEvaluator {
	e.keyGetter = keyGetter
	return e
}

func (e *BaseVarEvaluator) WithIgnoreComments(ignoreComments bool) *BaseVarEvaluator {
	e.ignoreComments = ignoreComments
	return e
}

func (e *BaseVarEvaluator) WithIgnoreInvalid(ignoreInvalid bool) *BaseVarEvaluator {
	e.ignoreInvalid = ignoreInvalid
	return e
}

func (e *BaseVarEvaluator) WithSkipRowColumnInfo(skipRowColumnInfo bool) *BaseVarEvaluator {
	e.skipRowColumnInfo = skipRowColumnInfo
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

func DefaultVarEncoder(c *VarContext, input interface{}) ([]byte, error) {
	return []byte("%v"), nil
}

func pathError(path []string, vars map[string]interface{}) error {
	keys := make([]string, 0, len(vars))

	for k := range vars {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	if len(path) == 0 {
		if len(keys) == 0 {
			return fmt.Errorf("no keys found")
		}

		return fmt.Errorf("possible keys are: %s", strings.Join(keys, ", "))
	}

	if len(keys) == 0 {
		return fmt.Errorf("no keys found for '%s'", strings.Join(path, "."))
	}

	return fmt.Errorf("possible keys for '%s' are: %s", strings.Join(path, "."), strings.Join(keys, ", "))
}

func DefaultVarKeyGetter(c *VarContext, vars map[string]interface{}) (val interface{}, err error) {
	var path []string

	parts := strings.Split(c.Token, ".")

	for _, part := range parts[:len(parts)-1] {
		varsint, ok := vars[part]
		if !ok {
			return nil, pathError(path, vars)
		}

		path = append(path, part)

		varsnext, ok := varsint.(map[string]interface{})
		if !ok {
			return nil, pathError(path, nil)
		}

		vars = varsnext
	}

	val, ok := vars[parts[len(parts)-1]]
	if !ok {
		return nil, pathError(path, vars)
	}

	return val, nil
}

func (e *BaseVarEvaluator) ExpandRaw(input []byte) (output []byte, params []interface{}, err error) {
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

		for start := range line {
			if start+1 >= ll || line[start] != e.varChar || line[start+1] != '{' {
				continue
			}

			idx := bytes.Index(line[start+2:], []byte("}"))
			if idx == -1 {
				continue
			}

			token = string(line[start+2 : start+2+idx])

			prefix := ""
			if !e.skipRowColumnInfo {
				prefix = fmt.Sprintf("[%d:%d] ", l+1, start+1)
			}

			if !validKey.MatchString(token) {
				if e.ignoreInvalid {
					continue
				}

				if token == "" {
					return nil, nil, fmt.Errorf("%sempty expansion found", prefix)
				}

				return nil, nil, fmt.Errorf("%sinvalid expansion found: %s", prefix, token)
			}

			out[l] = append(out[l], line[done:start]...)

			varCtx := &VarContext{
				Input:            input,
				Line:             line,
				Row:              l + 1,
				Token:            token,
				TokenColumnStart: start,
				TokenColumnEnd:   start + 2 + idx,
			}

			val, err := e.keyGetter(varCtx, e.vars)
			if err != nil {
				return nil, nil, fmt.Errorf("%sexpansion value for '%s' could not be evaluated:\n%w", prefix, token, err)
			}

			valOut, err := e.encoder(varCtx, val)
			if err != nil {
				return nil, nil, fmt.Errorf("%sexpansion value for '%s' could not be encoded, unknown field\nerror: %w", prefix, token, err)
			}

			out[l] = append(out[l], valOut...)
			params = append(params, val)

			done = start + 3 + idx
		}

		if done < ll {
			out[l] = append(out[l], line[done:]...)
		}
	}

	return bytes.Join(out, []byte{'\n'}), params, nil
}
