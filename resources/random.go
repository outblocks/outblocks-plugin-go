package resources

import (
	"context"
	"fmt"

	"github.com/outblocks/outblocks-plugin-go/registry"
	"github.com/outblocks/outblocks-plugin-go/registry/fields"
	"github.com/outblocks/outblocks-plugin-go/util"
)

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

func (o *RandomString) GetID() string {
	return fmt.Sprintf("randomstring/%s", o.Name.Any())
}

func (o *RandomString) GetName() string {
	return o.Name.Any()
}

func (o *RandomString) Create(ctx context.Context, meta any) error {
	res := util.RandomStringCryptoCustom(o.Lower.Wanted(), o.Upper.Wanted(), o.Numeric.Wanted(), o.Special.Wanted(), o.Length.Wanted())

	o.Result.SetCurrent(res)

	return nil
}

func (o *RandomString) Update(ctx context.Context, meta any) error {
	return fmt.Errorf("unimplemented")
}

func (o *RandomString) Delete(ctx context.Context, meta any) error {
	return nil
}
