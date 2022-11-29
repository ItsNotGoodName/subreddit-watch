package matcher

import (
	"regexp"

	"github.com/turnage/graw/reddit"
)

type Regexp struct {
	titleRegexp *regexp.Regexp
}

func NewRegexp(titleRegexp *regexp.Regexp) *Regexp {
	return &Regexp{
		titleRegexp: titleRegexp,
	}
}

func (r *Regexp) Match(p *reddit.Post) bool {
	return r.titleRegexp.Match([]byte(p.Title))
}
