package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type StringCommand struct {
	valStr string
	valArr []string
}

func (c *StringCommand) Array() []string {
	if c == nil {
		return nil
	}

	if len(c.valArr) != 0 {
		return c.valArr
	}

	if c.valStr != "" {
		return []string{c.valStr}
	}

	return nil
}

func (c *StringCommand) ArrayOrShell() []string {
	if c == nil {
		return nil
	}

	if c.IsArray() {
		return c.Array()
	}

	return []string{"sh", "-c", c.valStr}
}

func (c *StringCommand) ExecCmdAsUser() *exec.Cmd {
	if len(c.valArr) != 0 {
		return exec.CommandContext(context.TODO(), c.valArr[0], c.valArr[1:]...) //nolint:gosec
	}

	if c.valStr == "" {
		panic("no command value specified")
	}

	return NewCmdAsUser(c.valStr)
}

func (c *StringCommand) MarshalJSON() ([]byte, error) {
	if !c.IsEmpty() && c.IsArray() {
		return json.Marshal(c.Array())
	}

	return json.Marshal(c.valStr)
}

func (c *StringCommand) UnmarshalJSON(b []byte) error {
	var out any

	err := json.Unmarshal(b, &out)
	if err != nil {
		return err
	}

	return c.fromInterface(out)
}

func (c *StringCommand) UnmarshalYAML(dec func(any) error) error {
	var out any

	err := dec(&out)
	if err != nil {
		return err
	}

	return c.fromInterface(out)
}

func (c *StringCommand) UnmarshalMapstructure(o any) error {
	return c.fromInterface(o)
}

func (c *StringCommand) fromInterface(o any) error {
	if o == nil {
		return nil
	}

	if v, ok := o.(string); ok {
		c.valStr = v

		return nil
	}

	if v, ok := o.([]any); ok {
		c.valArr = make([]string, len(v))

		for i, val := range v {
			if s, ok := val.(string); ok {
				c.valArr[i] = s
			} else {
				return fmt.Errorf("invalid value, only string or array of strings are allowed")
			}
		}

		return nil
	}

	if v, ok := o.([]string); ok {
		c.valArr = v

		return nil
	}

	return fmt.Errorf("invalid value, only string or array of strings are allowed")
}

func (c *StringCommand) IsEmpty() bool {
	return c == nil || (len(c.valArr) == 0 && c.valStr == "")
}

func (c *StringCommand) IsArray() bool {
	return c != nil && len(c.valArr) != 0
}

func (c *StringCommand) Flatten() string {
	if c.valStr != "" {
		return c.valStr
	}

	ret := make([]string, len(c.valArr))
	for i, v := range c.valArr {
		ret[i] = strings.ReplaceAll(v, "\"", "\\\"")
	}

	return fmt.Sprintf("\"%s\"", strings.Join(ret, "\"")) //nolint:gocritic
}

func NewStringCommandFromString(s string) *StringCommand {
	return &StringCommand{
		valStr: s,
	}
}

func NewStringCommandFromArray(s []string) *StringCommand {
	return &StringCommand{
		valArr: s,
	}
}
