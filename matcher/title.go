package matcher

import (
	"regexp"

	"github.com/turnage/graw/reddit"
)

type TitleRegex struct {
	regexp *regexp.Regexp
}

func NewTitleRegex(title *regexp.Regexp) *TitleRegex {
	return &TitleRegex{
		regexp: title,
	}
}

func (tr *TitleRegex) Match(p *reddit.Post) bool {
	return tr.regexp.Match([]byte(p.Title))
}
