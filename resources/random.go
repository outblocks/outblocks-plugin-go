package resources

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/outblocks/outblocks-plugin-go/registry"
	"github.com/outblocks/outblocks-plugin-go/registry/fields"
)

const numChars = "0123456789"
const lowerChars = "abcdefghijklmnopqrstuvwxyz"
const upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const specialChars = "!@#$%&*()-_=+[]{}<>:?"

type RandomString struct {
	registry.ResourceBase

	Name    fields.StringInputField
	Length  fields.IntInputField  `state:"force_new" default:"16"`
	Lower   fields.BoolInputField `state:"force_new" default:"true"`
	Upper   fields.BoolInputField `state:"force_new" default:"true"`
	Numeric fields.BoolInputField `state:"force_new" default:"true"`
	Special fields.BoolInputField `state:"force_new" default:"true"`

	Result fields.StringOutputField
}

func (o *RandomString) GetName() string {
	return o.Name.Any()
}

func (o *RandomString) Create(ctx context.Context, meta interface{}) error {
	var chars string

	if o.Lower.Wanted() {
		chars += lowerChars
	}

	if o.Upper.Wanted() {
		chars += upperChars
	}

	if o.Numeric.Wanted() {
		chars += numChars
	}

	if o.Special.Wanted() {
		chars += specialChars
	}

	bytes := make([]byte, o.Length.Wanted())
	setLen := big.NewInt(int64(len(chars)))

	for i := range bytes {
		idx, err := rand.Int(rand.Reader, setLen)
		if err != nil {
			return err
		}

		bytes[i] = chars[idx.Int64()]
	}

	o.Result.SetCurrent(string(bytes))

	return nil
}

func (o *RandomString) Update(ctx context.Context, meta interface{}) error {
	return fmt.Errorf("unimplemented")
}

func (o *RandomString) Delete(ctx context.Context, meta interface{}) error {
	return nil
}
