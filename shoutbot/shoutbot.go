package shoutbot

import (
	"fmt"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/turnage/graw/reddit"
)

type Templater interface {
	GetTitle(p *reddit.Post) string
	GetMessage(p *reddit.Post) string
}

type Matcher interface {
	Match(p *reddit.Post) bool
	String() string
}

type ShoutBot struct {
	paths map[string]Path
}

type Path struct {
	Matcher   Matcher
	Templater Templater
	Sender    *router.ServiceRouter
}

func New(paths map[string]Path) *ShoutBot {
	return &ShoutBot{
		paths: paths,
	}
}

func (sb *ShoutBot) Post(p *reddit.Post) error {
	entry := NewLogEntry(p)
	defer func() { fmt.Println(entry.Trace()) }()

	// Get path by subreddit
	path, ok := sb.paths[p.Subreddit]
	if !ok {
		entry.Log("ERROR: invalid subreddit received: " + p.Subreddit)
		return nil
	}

	// Check if post matches
	if ok = path.Matcher.Match(p); !ok {
		entry.Log("FAIL: " + path.Matcher.String())
		return nil
	}
	entry.Log("PASS: " + path.Matcher.String())

	// Send message to sender
	r := types.Params{}
	r.SetTitle(path.Templater.GetTitle(p))
	errs := path.Sender.Send(path.Templater.GetMessage(p), &r)
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	entry.Log("PASS: sent notification")

	return nil
}
