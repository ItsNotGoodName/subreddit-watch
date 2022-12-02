package matcher

import (
	"regexp"

	"github.com/turnage/graw/reddit"
)

type TitleRegex struct {
	regexp *regexp.Regexp
	regstr string
}

func NewTitleRegex(title *regexp.Regexp) *TitleRegex {
	return &TitleRegex{
		regexp: title,
		regstr: title.String(),
	}
}

func (tr *TitleRegex) Match(p *reddit.Post) bool {
	return tr.regexp.MatchString(p.Title)
}

func (tr *TitleRegex) String() string {
	return "title regex matcher: " + tr.regstr
}
