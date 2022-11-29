package matcher

import (
	"github.com/turnage/graw/reddit"
)

type Default struct{}

func NewDefault() *Default {
	return &Default{}
}

func (Default) Match(p *reddit.Post) bool {
	return true
}
