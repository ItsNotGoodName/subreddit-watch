# subreddit-watch

Watch subreddits for new posts.

# Configuration

Program looks for `./config.yml` or `~/.subreddit-watch.yml`.

```yaml
reddit_id: XXXXXXXXXXXXXXXXXXXXXX
reddit_secret: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
redit_username: XXXXXXXXXXXXXXXXXXXX # Optional

notify_title_template: {{ .Post.Title }}  # Optional
notify_message_template: |                # Optional
  https://old.reddit.com{{ .Post.Permalink }}{{ if not .Post.IsSelf }}
  {{ .Post.URL }}{{ end }}
notify: # https://containrrr.dev/shoutrrr/v0.5/services/overview/
  - telegram://XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX@telegram?chats=-1001111111111&Preview=false

subreddits:
  - name: buildapcsales
    title_regex: (?i)^\[GPU\]                                 # Optional
    notify_title_template: buildapcsales - {{ .Post.Title }}  # Optional
    notify_message_template: |                                # Optional
      https://old.reddit.com{{ .Post.Permalink }}{{ if not .Post.IsSelf }}
      {{ .Post.URL }}{{ end }}
```