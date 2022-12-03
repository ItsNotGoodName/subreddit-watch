package config

import (
	"fmt"

	"github.com/ItsNotGoodName/subreddit-watch/templater"
	"github.com/spf13/viper"
)

type Raw struct {
	RedditSecret   string `mapstructure:"reddit_secret"`
	RedditID       string `mapstructure:"reddit_id"`
	RedditUsername string `mapstructure:"reddit_username"`

	NotifyTitleTemplate   string   `mapstructure:"notify_title_template"`
	NotifyMessageTemplate string   `mapstructure:"notify_message_template"`
	Notify                []string `mapstructure:"notify"`

	Subreddits []struct {
		Name                  string   `mapstructure:"name"`
		TitleRegex            []string `mapstructure:"title_regex"`
		NotifyTitleTemplate   string   `mapstructure:"notify_title_template"`
		NotifyMessageTemplate string   `mapstructure:"notify_message_template"`
	} `mapstructure:"subreddits"`
}

func Read(configFile string) (*Raw, error) {
	viper.SetConfigName("config")                 // name of config file (without extension)
	viper.SetConfigType("yaml")                   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/subreddit-watch/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.subreddit-watch") // call multiple times to add many search paths
	viper.AddConfigPath(".")                      // optionally look for config in the working directory
	viper.SetDefault("notify_title_template", templater.DefaultTitleTemplate)
	viper.SetDefault("notify_message_template", templater.DefaultMessageTemplate)
	viper.RegisterAlias("subreddit", "subreddits")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	raw := &Raw{}
	if err := viper.Unmarshal(raw); err != nil {
		return nil, err
	}

	// Prevent duplicate subreddit
	subreddits := make(map[string]struct{})
	for _, v := range raw.Subreddits {
		if _, ok := subreddits[v.Name]; ok {
			return nil, fmt.Errorf("duplicate subreddit: %s", v.Name)
		}
		subreddits[v.Name] = struct{}{}
	}

	return raw, nil
}
