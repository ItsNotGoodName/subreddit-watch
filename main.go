package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/ItsNotGoodName/subreddit-watch/config"
	"github.com/ItsNotGoodName/subreddit-watch/matcher"
	"github.com/ItsNotGoodName/subreddit-watch/shoutbot"
	"github.com/ItsNotGoodName/subreddit-watch/templater"
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

const AppID string = "github.com/ItsNotGoodName/subreddit-watch"

var (
	Platform string = runtime.GOOS
	Version  string = "dev"
)

func main() {
	argTest := flag.Bool("test", false, "test current configuration")
	argVersion := flag.Bool("version", false, "show version")

	flag.Parse()

	if *argVersion {
		fmt.Println(Version)
		return
	}

	// Parse config
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalln("main: config parse error:", err)
	}

	// Create sender
	sender, err := shoutrrr.CreateSender(cfg.Notify...)
	if err != nil {
		log.Fatalln("main: sender create error:", err)
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

	runner := shoutbot.New(newPaths(cfg, sender))

	if *argTest {
		// Test command
		fmt.Println("Testing")

		if len(cfg.Subreddits) == 0 {
			log.Fatalln("main: no subreddits defined")
		}

		harvest, err := bot.ListingWithParams("/r/"+cfg.Subreddits[0].Name+"/new", map[string]string{"limit": "1"})
		if err != nil {
			log.Fatalln("main: list error:", err)
		}

		if len(harvest.Posts) == 0 {
			log.Fatalln("main: invalid posts length:", len(harvest.Posts))
		}

		if err := runner.Post(harvest.Posts[0]); err != nil {
			log.Fatalln("main: bot post error:", err)
		}
	} else {
		// Run command
		if _, wait, err := graw.Run(runner, bot, graw.Config{Subreddits: cfg.SubredditNameList()}); err != nil {
			log.Println("main: graw run error:", err)
		} else {
			fmt.Println("Started")
			log.Println("main: graw wait error:", wait())
		}
	}
}

func newPaths(c *config.Config, sender *router.ServiceRouter) map[string]shoutbot.Path {
	paths := make(map[string]shoutbot.Path)
	for _, subreddit := range c.Subreddits {
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
	}
	return paths
}