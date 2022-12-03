package shoutbot

import (
	"fmt"

	"github.com/turnage/graw/reddit"
)

type TraceEntry struct {
	id   string
	logs []string
}

func NewTraceEntry(p *reddit.Post) *TraceEntry {
	return &TraceEntry{
		id:   p.ID,
		logs: []string{fmt.Sprintf("%s: %s", p.Title, p.Permalink)},
	}
}

func (te *TraceEntry) Pass(log string) {
	te.logs = append(te.logs, "PASS: "+log)
}

func (te *TraceEntry) Fail(log string) {
	te.logs = append(te.logs, "FAIL: "+log)
}

func (te *TraceEntry) Error(log string) {
	te.logs = append(te.logs, "ERROR: "+log)
}

func (te *TraceEntry) Trace() string {
	var ret string
	for _, log := range te.logs {
		ret += te.id + ": " + log + "\n"
	}
	return ret
}
