package types

import (
	"strings"

	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
)

type DomainInfoMatcher struct {
	normal    map[string]*apiv1.DomainInfo
	wildcards map[string]*apiv1.DomainInfo
}

func NewDomainInfoMatcher(domains []*apiv1.DomainInfo) *DomainInfoMatcher {
	m := &DomainInfoMatcher{
		normal:    make(map[string]*apiv1.DomainInfo),
		wildcards: make(map[string]*apiv1.DomainInfo),
	}

	for _, o := range domains {
		for _, d := range o.Domains {
			dparts := strings.SplitN(d, ".", 2)
			if len(dparts) == 2 && dparts[0] == "*" {
				m.wildcards[dparts[1]] = o
			} else {
				m.normal[d] = o
			}
		}
	}

	return m
}

func (m *DomainInfoMatcher) Match(v string) *apiv1.DomainInfo {
	di := m.normal[v]
	if di != nil {
		return di
	}

	dparts := strings.SplitN(v, ".", 2)

	di = m.wildcards[dparts[1]]
	if di != nil {
		return di
	}

	return nil
}
