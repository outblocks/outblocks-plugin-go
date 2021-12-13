package types

import (
	"regexp"

	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
	"github.com/outblocks/outblocks-plugin-go/util"
)

type DomainInfoMatcher struct {
	matchers map[*regexp.Regexp]*apiv1.DomainInfo
}

func NewDomainInfoMatcher(domains []*apiv1.DomainInfo) *DomainInfoMatcher {
	m := &DomainInfoMatcher{
		matchers: make(map[*regexp.Regexp]*apiv1.DomainInfo),
	}

	for _, o := range domains {
		for _, d := range o.Domains {
			re := util.DomainRegex(d)
			m.matchers[re] = o
		}
	}

	return m
}

func (m *DomainInfoMatcher) Match(v string) *apiv1.DomainInfo {
	for m, d := range m.matchers {
		if m.MatchString(v) {
			return d
		}
	}

	return nil
}
