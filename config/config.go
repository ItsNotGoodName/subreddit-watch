package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"text/template"

	"github.com/ItsNotGoodName/subreddit-watch/templater"
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/spf13/viper"
)

type Raw struct {
	RedditSecret   string `mapstructure:"reddit_secret"`
	RedditID       string `mapstructure:"reddit_id"`
	RedditUsername string `mapstructure:"reddit_username"`

	NotifyTitleTemplate   string   `mapstructure:"notify_title_template"`
	NotifyMessageTemplate string   `mapstructure:"notify_message_template"`
	Notify                []string `mapstructure:"notify"`

	Subreddits []RawSubreddits `mapstructure:"subreddits"`
}

type RawSubreddits struct {
	Name                  string   `mapstructure:"name"`
	TitleRegex            []string `mapstructure:"title_regex"`
	NotifyTitleTemplate   string   `mapstructure:"notify_title_template"`
	NotifyMessageTemplate string   `mapstructure:"notify_message_template"`
	Notify                []string `mapstructure:"notify"`
}

func InitConfig(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln("config.InitConfig:", err)
		}

		// Search config in home directory with name ".subreddit-watch" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".subreddit-watch")
	}

	viper.AutomaticEnv() // read in environment variables that match

	viper.SetDefault("notify_title_template", templater.DefaultTitleTemplate)
	viper.SetDefault("notify_message_template", templater.DefaultMessageTemplate)

	viper.RegisterAlias("subreddit", "subreddits")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("config.InitConfig:", err)
	}
}

func Parse() (*Config, error) {
	raw := &Raw{}
	if err := viper.Unmarshal(raw); err != nil {
		return nil, err
	}

	// Need atleast 1 subreddit to run
	if raw.Subreddits == nil {
		return nil, fmt.Errorf("no subreddits defined")
	}

	// Prevent duplicate subreddit
	subreddits := make(map[string]struct{})
	for _, v := range raw.Subreddits {
		if _, ok := subreddits[v.Name]; ok {
			return nil, fmt.Errorf("duplicate subreddit: %s", v.Name)
		}
		subreddits[v.Name] = struct{}{}
	}

	return parse(raw)
}

type Config struct {
	RedditSecret   string
	RedditID       string
	RedditUsername string

	Subreddits []Subreddits
}

func (c *Config) SubredditNameList() []string {
	subreddits := make([]string, len(c.Subreddits))
	for i := range c.Subreddits {
		subreddits[i] = c.Subreddits[i].Name
	}
	return subreddits
}

type Subreddits struct {
	Name                  string
	TitleRegex            []*regexp.Regexp
	NotifyTitleTemplate   *template.Template
	NotifyMessageTemplate *template.Template
	Notify                *router.ServiceRouter
}

func parse(raw *Raw) (*Config, error) {
	defaultNotifyTitleTemplate, err := template.New("").Parse(raw.NotifyTitleTemplate)
	if err != nil {
		return nil, err
	}

	defaultNotifyMessageTemplate, err := template.New("").Parse(raw.NotifyMessageTemplate)
	if err != nil {
		return nil, err
	}

	defaultNotify, err := shoutrrr.CreateSender(raw.Notify...)
	if err != nil {
		return nil, err
	}

	titleRegexMap := make(map[string]*regexp.Regexp)
	subreddits := make([]Subreddits, len(raw.Subreddits))
	for i, s := range raw.Subreddits {
		// Title regex
		titleRegex := make([]*regexp.Regexp, len(s.TitleRegex))
		for i, expr := range s.TitleRegex {
			if reg, ok := titleRegexMap[expr]; ok {
				titleRegex[i] = reg
			} else {
				reg, err := regexp.Compile(expr)
				if err != nil {
					return nil, err
				}

				titleRegex[i] = reg
				titleRegexMap[expr] = reg
			}
		}

		// Title template
		notifyTitleTemplate := defaultNotifyTitleTemplate
		if raw.NotifyTitleTemplate != s.NotifyTitleTemplate {
			var err error
			notifyTitleTemplate, err = template.New("").Parse(raw.NotifyTitleTemplate)
			if err != nil {
				return nil, err
			}
		}

		// Message template
		notifyMessageTemplate := defaultNotifyMessageTemplate
		if raw.NotifyMessageTemplate != s.NotifyMessageTemplate {
			var err error
			notifyMessageTemplate, err = template.New("").Parse(raw.NotifyMessageTemplate)
			if err != nil {
				return nil, err
			}
		}

		// Notify
		notify := defaultNotify
		if s.Notify != nil {
			var err error
			notify, err = shoutrrr.CreateSender(s.Notify...)
			if err != nil {
				return nil, err
			}
		}

		subreddits[i] = Subreddits{
			Name:                  s.Name,
			TitleRegex:            titleRegex,
			NotifyTitleTemplate:   notifyTitleTemplate,
			NotifyMessageTemplate: notifyMessageTemplate,
			Notify:                notify,
		}
	}

	return &Config{
		RedditSecret:   raw.RedditSecret,
		RedditID:       raw.RedditID,
		RedditUsername: raw.RedditUsername,
		Subreddits:     subreddits,
	}, nil
}
