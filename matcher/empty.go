package matcher

import (
	"github.com/turnage/graw/reddit"
)

type Empty struct{}

func NewEmpty() *Empty {
	return &Empty{}
}

func (Empty) Match(p *reddit.Post) bool {
	return true
}
