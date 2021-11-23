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
