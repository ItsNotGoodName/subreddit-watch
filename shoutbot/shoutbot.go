package shoutbot

import (
	"fmt"
	"log"

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
	fmt.Printf("%s: %s: %s\n", p.ID, p.Title, p.Permalink)

	// Get path by subreddit
	path, ok := sb.paths[p.Subreddit]
	if !ok {
		log.Printf("shoutbot.Shoutbot.Post: subreddit path not found error: %s", p.Subreddit)
		return nil
	}

	// Check if post matches
	if path.Matcher.Match(p); !ok {
		return nil
	}

	// Send message to sender
	r := types.Params{}
	r.SetTitle(path.Templater.GetTitle(p))
	errs := path.Sender.Send(path.Templater.GetMessage(p), &r)
	for _, e := range errs {
		return e
	}

	return nil
}
