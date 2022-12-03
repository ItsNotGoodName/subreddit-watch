package config

import (
	"regexp"
	"text/template"
)

type Config struct {
	RedditSecret   string
	RedditID       string
	RedditUsername string

	Notify []string

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
}

func Parse(raw *Raw) (*Config, error) {
	defaultNotifyTitleTemplate, err := template.New("").Parse(raw.NotifyTitleTemplate)
	if err != nil {
		return nil, err
	}

	defaultNotifyMessageTemplate, err := template.New("").Parse(raw.NotifyMessageTemplate)
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

		subreddits[i] = Subreddits{
			Name:                  s.Name,
			TitleRegex:            titleRegex,
			NotifyTitleTemplate:   notifyTitleTemplate,
			NotifyMessageTemplate: notifyMessageTemplate,
		}
	}

	return &Config{
		RedditSecret:   raw.RedditSecret,
		RedditID:       raw.RedditID,
		RedditUsername: raw.RedditUsername,
		Notify:         raw.Notify,
		Subreddits:     subreddits,
	}, nil
}