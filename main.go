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
	argConfig := flag.String("config", "", "config file")
	argVersion := flag.Bool("version", false, "show version")

	flag.Parse()

	if *argVersion {
		fmt.Println(Version)
		return
	}

	// Read and parse config
	cfg, err := config.ParseAfter(config.Read(*argConfig))
	if err != nil {
		log.Fatalln("main: config parse error:", err)
	}

	// Create sender
	sender, err := shoutrrr.CreateSender(cfg.Notify...)
	if err != nil {
		log.Fatalln("main: sender create error:", err)
	}

	// Create bot
	botConfig := reddit.BotConfig{
		Agent: fmt.Sprintf("%s:%s:%s (by /u/%s)", Platform, AppID, Version, cfg.RedditUsername),
		App: reddit.App{
			ID:     cfg.RedditID,
			Secret: cfg.RedditSecret,
		},
	}
	bot, err := reddit.NewBot(botConfig)
	if err != nil {
		log.Fatalln("main: bot create error:", err)
	}

	handler := shoutbot.New(newPaths(cfg, sender))

	if *argTest {
		// Test command
		if len(cfg.Subreddits) == 0 {
			log.Fatalln("main: no subreddits defined")
		}
		fmt.Printf("Started test for '%s'\n", cfg.Subreddits[0].Name)

		harvest, err := bot.ListingWithParams("/r/"+cfg.Subreddits[0].Name+"/new", map[string]string{"limit": "1"})
		if err != nil {
			log.Fatalln("main: list error:", err)
		}

		if len(harvest.Posts) == 0 {
			log.Fatalln("main: invalid posts length:", len(harvest.Posts))
		}

		if err := handler.Post(harvest.Posts[0]); err != nil {
			log.Fatalln("main: bot post error:", err)
		}
	} else {
		// Run command
		subreddits := cfg.SubredditNameList()
		if _, wait, err := graw.Run(handler, bot, graw.Config{Subreddits: subreddits}); err != nil {
			log.Println("main: graw run error:", err)
		} else {
			fmt.Println("Watching", subreddits)
			log.Println("main: graw wait error:", wait())
		}
	}
}

func newPaths(c *config.Config, sender *router.ServiceRouter) map[string]shoutbot.Path {
	paths := make(map[string]shoutbot.Path)
	for _, subreddit := range c.Subreddits {
		matchers := []shoutbot.Matcher{}
		for _, reg := range subreddit.TitleRegex {
			matchers = append(matchers, matcher.NewTitleRegex(reg))
		}

		paths[subreddit.Name] = shoutbot.Path{
			Sender:    sender,
			Matchers:  matchers,
			Templater: templater.New(subreddit.NotifyTitleTemplate, subreddit.NotifyMessageTemplate),
		}
	}
	return paths
}
