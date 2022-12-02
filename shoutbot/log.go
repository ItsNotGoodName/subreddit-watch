package shoutbot

import (
	"fmt"

	"github.com/turnage/graw/reddit"
)

type LogEntry struct {
	id   string
	logs []string
}

func NewLogEntry(p *reddit.Post) *LogEntry {
	return &LogEntry{
		id:   p.ID,
		logs: []string{fmt.Sprintf("%s: %s", p.Title, p.Permalink)},
	}
}

func (le *LogEntry) Log(log string) {
	le.logs = append(le.logs, log)
}

func (le *LogEntry) Trace() string {
	var ret string
	for _, log := range le.logs {
		ret += le.id + ": " + log + "\n"
	}
	return ret
}
