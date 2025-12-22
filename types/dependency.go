package types

import (
	"fmt"
	"strings"

	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
	"github.com/outblocks/outblocks-plugin-go/util"
)

func DependencyEnvPrefix(d *apiv1.Dependency) string {
	return fmt.Sprintf("DEP_%s_%s_", strings.ToUpper(d.Type), util.SanitizeEnvVar((strings.ToUpper(d.Name))))
}

type StorageDepOptions struct {
	Name               string `json:"name"`
	Versioning         bool   `json:"versioning"`
	Location           string `json:"location"`
	DeleteInDays       int    `json:"delete_in_days"`
	ExpireVersionsDays int    `json:"expire_versions_in_days"`
	MaxVersions        int    `json:"max_versions"`
	Public             bool   `json:"public"`

	CORS []struct {
		Origins         []string `json:"origins"`
		Methods         []string `json:"methods"`
		ResponseHeaders []string `json:"response_headers"`
		MaxAgeInSeconds int      `json:"max_age_in_seconds"`
	} `json:"cors"`
}

func NewStorageDepOptions(in map[string]any) (*StorageDepOptions, error) {
	o := &StorageDepOptions{}

	return o, util.MapstructureJSONDecode(in, o)
}

type DatabaseDepOptionUser struct {
	Password string `json:"password"`
	Hostname string `json:"hostname"`
}

type DatabaseDepOptions struct {
	Version string                            `json:"version"`
	HA      bool                              `json:"high_availability"`
	Tier    string                            `json:"tier"`
	Flags   map[string]string                 `json:"flags"`
	Users   map[string]*DatabaseDepOptionUser `json:"users"`
}

func NewDatabaseDepOptions(in map[string]any) (*DatabaseDepOptions, error) {
	o := &DatabaseDepOptions{}

	return o, util.MapstructureJSONDecode(in, o)
}

type DatabaseDepNeed struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Database string `json:"database"`
}

func NewDatabaseDepNeed(in map[string]any) (*DatabaseDepNeed, error) {
	o := &DatabaseDepNeed{}

	return o, util.MapstructureJSONDecode(in, o)
}
