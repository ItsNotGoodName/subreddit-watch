FROM alpine
ENTRYPOINT ["/usr/bin/subreddit-watch"]
WORKDIR /config
COPY subreddit-watch /usr/bin/subreddit-watch
