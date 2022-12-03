# subreddit-watch

Watch subreddits for new posts and send notifications.

# Configuration

Configuration is loaded from `./config.yml`, `~/.subreddit-watch.yml`, or `/etc/subreddit-watch.yml`.

```yaml
# Reddit API access
reddit_id: XXXXXXXXXXXXXXXXXXXXXX
reddit_secret: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
redit_username: XXXXXXXXXXXXXXXXXXXX # Optional

# Default notifications
notify_title_template: {{ .Post.Title }}  # Optional
notify_message_template: |                # Optional
  https://old.reddit.com{{ .Post.Permalink }}{{ if not .Post.IsSelf }}
  {{ .Post.URL }}{{ end }}
notify: # https://containrrr.dev/shoutrrr/v0.5/services/overview/
  - telegram://XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX@telegram?chats=-1001111111111&Preview=false

# Subreddits to watch
subreddits:
  - name: buildapcsales
    title_regex:                              # Optional
      - (?i)^\[GPU\]
    notify_title_template: {{ .Post.Title }}  # Optional
    notify_message_template: |                # Optional
      https://old.reddit.com{{ .Post.Permalink }}{{ if not .Post.IsSelf }}
      {{ .Post.URL }}{{ end }}
```