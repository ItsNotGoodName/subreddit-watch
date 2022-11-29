package config

import (
	"fmt"
	"regexp"
	"text/template"

	"github.com/ItsNotGoodName/subreddit-watch/templater"
	"github.com/spf13/viper"
)

type Config struct {
	RedditSecret   string `mapstructure:"reddit_secret"`
	RedditID       string `mapstructure:"reddit_id"`
	RedditUsername string `mapstructure:"reddit_username"`

	NotifyTitleTemplate      *template.Template `mapstructure:"-"`
	NotifyTitleTemplateRaw   string             `mapstructure:"notify_title_template"`
	NotifyMessageTemplate    *template.Template `mapstructure:"-"`
	NotifyMessageTemplateRaw string             `mapstructure:"notify_message_template"`
	Notify                   []string           `mapstructure:"notify"`

	Subreddits []*Subreddit `mapstructure:"subreddits"`
}

type Subreddit struct {
	Name                     string             `mapstructure:"name"`
	TitleRegex               *regexp.Regexp     `mapstructure:"-"`
	TitleRegexRaw            string             `mapstructure:"title_regex"`
	NotifyTitleTemplate      *template.Template `mapstructure:"-"`
	NotifyTitleTemplateRaw   string             `mapstructure:"notify_title_template"`
	NotifyMessageTemplate    *template.Template `mapstructure:"-"`
	NotifyMessageTemplateRaw string             `mapstructure:"notify_message_template"`
}

func Parse() (*Config, error) {
	viper.SetConfigName("config")                 // name of config file (without extension)
	viper.SetConfigType("yaml")                   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/subreddit-watch/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.subreddit-watch") // call multiple times to add many search paths
	viper.AddConfigPath(".")                      // optionally look for config in the working directory
	viper.SetDefault("notify_title_template", templater.DefaultTitleTemplate)
	viper.SetDefault("notify_message_template", templater.DefaultMessageTemplate)
	viper.RegisterAlias("subreddit", "subreddits")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Notify Title and Message Template
	cfg.NotifyTitleTemplate = template.Must(template.New("").Parse(cfg.NotifyTitleTemplateRaw))
	cfg.NotifyMessageTemplate = template.Must(template.New("").Parse(cfg.NotifyMessageTemplateRaw))

	subreddits := make(map[string]struct{})
	for _, v := range cfg.Subreddits {
		// Prevent duplicate subreddit
		if _, ok := subreddits[v.Name]; ok {
			return nil, fmt.Errorf("duplicate subreddit: %s", v.Name)
		}
		subreddits[v.Name] = struct{}{}

		if v.TitleRegexRaw != "" {
			v.TitleRegex = regexp.MustCompile(v.TitleRegexRaw)
		}

		// Notify Title Template
		if v.NotifyTitleTemplateRaw == "" {
			v.NotifyTitleTemplate = cfg.NotifyTitleTemplate
		} else {
			v.NotifyTitleTemplate = template.Must(template.New("").Parse(v.NotifyTitleTemplateRaw))
		}

		// Notify Message Template
		if v.NotifyMessageTemplateRaw == "" {
			v.NotifyMessageTemplate = cfg.NotifyMessageTemplate
		} else {
			v.NotifyMessageTemplate = template.Must(template.New("").Parse(v.NotifyMessageTemplateRaw))
		}
	}

	return cfg, nil
}
