# subreddit-watch

Watch and notify when subreddits have new posts.

# Configuration

Configuration is loaded from `./.subreddit-watch.yml`, `~/.subreddit-watch.yml`, or `/etc/.subreddit-watch.yml`.

Each template has access to the [`reddit.Post`](https://pkg.go.dev/github.com/turnage/graw/reddit#Post) variable via the `.Post` template variable.

```yaml
# Reddit API access https://www.reddit.com/prefs/apps
reddit_id: XXXXXXXXXXXXXXXXXXXXXX
reddit_secret: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
redit_username: XXXXXXXXXXXXXXXXXXXX # Optional

# Default notifications
notify_title_template: "{{ .Post.Title }}" # Optional
notify_message_template: | # Optional
  https://old.reddit.com{{ .Post.Permalink }}{{ if not .Post.IsSelf }}
  {{ .Post.URL }}{{ end }}
notify: # Optional https://containrrr.dev/shoutrrr/v0.7/services/overview/
  - telegram://XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX@telegram?chats=-1001111111111&Preview=false

# Subreddits to watch
subreddits:
  - name: buildapcsales
    title_regex: # Optional
      - (?i)^\[GPU\]
    notify_title_template: "{{ .Post.Title }}" # Optional
    notify_message_template: | # Optional
      https://old.reddit.com{{ .Post.Permalink }}{{ if not .Post.IsSelf }}
      {{ .Post.URL }}{{ end }}
    notify: # Optional https://containrrr.dev/shoutrrr/v0.7/services/overview/
      - telegram://XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX@telegram?chats=-1001111111111&Preview=false
```

# Docker

Configuration is loaded from `/config/.subreddit-watch.yml`.

## docker-compose

```yaml
version: "3"
services:
  subreddit-watch:
    container_name: subreddit-watch
    image: ghcr.io/itsnotgoodname/subreddit-watch:latest
    volumes:
      - /path/to/appdata/config:/config
    user: 1000:1000
    restart: unless-stopped
```

## docker cli

```shell
docker run -d \
  --name=subreddit-watch \
  -v /path/to/appdata/config:/config \
  --user 1000:1000 \
  --restart unless-stopped \
  ghcr.io/itsnotgoodname/subreddit-watch:latest
```
