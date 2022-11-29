package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/ItsNotGoodName/subreddit-watch/config"
	"github.com/ItsNotGoodName/subreddit-watch/matcher"
	"github.com/ItsNotGoodName/subreddit-watch/shoutbot"
	"github.com/ItsNotGoodName/subreddit-watch/templater"
	"github.com/containrrr/shoutrrr"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

const AppID string = "github.com/ItsNotGoodName/subreddit-watch"

var (
	Platform string = runtime.GOOS
	Version  string = "dev"
)

func main() {
	// Parse config
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalln("main: config parse error:", err)
	}

	// Create bot
	bcfg := reddit.BotConfig{
		Agent: fmt.Sprintf("%s:%s:%s (by /u/%s)", Platform, AppID, Version, cfg.RedditUsername),
		App: reddit.App{
			ID:     cfg.RedditID,
			Secret: cfg.RedditSecret,
		},
	}
	bot, err := reddit.NewBot(bcfg)
	if err != nil {
		log.Fatalln("main: bot create error:", err)
	}

	// Create sender
	sender, err := shoutrrr.CreateSender(cfg.Notify...)
	if err != nil {
		log.Println("main: sender create error:", err)
	}

	paths := make(map[string]shoutbot.Path)
	subreddits := []string{}
	for _, subreddit := range cfg.Subreddits {
		var m shoutbot.Matcher
		if subreddit.TitleRegex != nil {
			m = matcher.NewTitleRegex(subreddit.TitleRegex)
		} else {
			m = matcher.NewDefault()
		}
		paths[subreddit.Name] = shoutbot.Path{
			Sender:    sender,
			Matcher:   m,
			Templater: templater.New(subreddit.NotifyTitleTemplate, subreddit.NotifyMessageTemplate),
		}
		subreddits = append(subreddits, subreddit.Name)
	}

	// Run bot
	gcfg := graw.Config{Subreddits: subreddits}
	if _, wait, err := graw.Run(shoutbot.NewShoutBot(paths), bot, gcfg); err != nil {
		log.Println("main: graw run error:", err)
	} else {
		fmt.Println("Started")
		log.Println("main: graw wait error:", wait())
	}
}
