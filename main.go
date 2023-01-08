package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/ItsNotGoodName/subreddit-watch/bot"
	"github.com/ItsNotGoodName/subreddit-watch/config"
	"github.com/ItsNotGoodName/subreddit-watch/matcher"
	"github.com/ItsNotGoodName/subreddit-watch/templater"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

var Version string = "dev"

func main() {
	argTest := flag.Bool("test", false, "test current configuration")
	argConfig := flag.String("config", "", "config file")
	argVersion := flag.Bool("version", false, "show version")

	flag.Parse()

	// Version
	if *argVersion {
		fmt.Println(Version)
		return
	}

	config.InitConfig(*argConfig)

	// Read and parse config
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalln("main: config parse error:", err)
	}

	// Create bot
	redditBot, err := reddit.NewBot(reddit.BotConfig{
		Agent: fmt.Sprintf("%s:%s:%s (by /u/%s)", runtime.GOOS, "github.com/ItsNotGoodName/subreddit-watch", Version, cfg.RedditUsername),
		App: reddit.App{
			ID:     cfg.RedditID,
			Secret: cfg.RedditSecret,
		},
	})
	if err != nil {
		log.Fatalln("main: bot create error:", err)
	}

	shoutBot := bot.NewShoutBot(newPaths(cfg))

	if *argTest {
		// Test command
		fmt.Printf("Started test for '%s'\n", cfg.Subreddits[0].Name)

		harvest, err := redditBot.ListingWithParams("/r/"+cfg.Subreddits[0].Name+"/new", map[string]string{"limit": "1"})
		if err != nil {
			log.Fatalln("main: list error:", err)
		}

		if len(harvest.Posts) == 0 {
			log.Fatalln("main: invalid posts length:", len(harvest.Posts))
		}

		if err := shoutBot.Post(harvest.Posts[0]); err != nil {
			log.Fatalln("main: bot post error:", err)
		}
	} else {
		// Run command
		subreddits := cfg.SubredditNameList()
		if _, wait, err := graw.Run(shoutBot, redditBot, graw.Config{Subreddits: subreddits}); err != nil {
			log.Println("main: graw run error:", err)
		} else {
			fmt.Println("Watching", subreddits)
			log.Println("main: graw wait error:", wait())
		}
	}
}

func newPaths(c *config.Config) map[string]bot.Path {
	paths := make(map[string]bot.Path)
	for _, subreddit := range c.Subreddits {
		matchers := []bot.Matcher{}
		for _, reg := range subreddit.TitleRegex {
			matchers = append(matchers, matcher.NewTitleRegex(reg))
		}

		paths[subreddit.Name] = bot.Path{
			Sender:    subreddit.Notify,
			Matchers:  matchers,
			Templater: templater.New(subreddit.NotifyTitleTemplate, subreddit.NotifyMessageTemplate),
			History:   bot.NewHistory(10),
		}
	}
	return paths
}
