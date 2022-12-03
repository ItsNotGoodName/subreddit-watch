package bot

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
	Matchers  []Matcher
	Templater Templater
	Sender    *router.ServiceRouter
}

func NewShoutBot(paths map[string]Path) *ShoutBot {
	return &ShoutBot{
		paths: paths,
	}
}

func (sb *ShoutBot) Post(p *reddit.Post) error {
	entry := NewTraceEntry(p)
	defer func() { fmt.Println(entry.Trace()) }()

	// Get path by subreddit
	path, ok := sb.paths[p.Subreddit]
	if !ok {
		entry.Error("invalid subreddit received: " + p.Subreddit)
		return nil
	}

	// Check if post matches
	if len(path.Matchers) > 0 {
		passed := false
		for _, matcher := range path.Matchers {
			if ok = matcher.Match(p); ok {
				entry.Pass(fmt.Sprintf("matcher %d: %s", len(path.Matchers), matcher.String()))
				passed = true
				break
			}
		}
		if !passed {
			entry.Fail(fmt.Sprintf("%d matchers", len(path.Matchers)))
			return nil
		}
	} else {
		entry.Pass(fmt.Sprintf("%d matchers", len(path.Matchers)))
	}

	// Send message to sender
	r := types.Params{}
	r.SetTitle(path.Templater.GetTitle(p))
	errs := path.Sender.Send(path.Templater.GetMessage(p), &r)
	for _, e := range errs {
		if e != nil {
			entry.Error("notification: " + e.Error())
			return nil
		}
	}
	entry.Pass("sent notification")

	return nil
}
