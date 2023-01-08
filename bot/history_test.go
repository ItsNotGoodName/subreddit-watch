package bot

import (
	"testing"

	"github.com/turnage/graw/reddit"
)

func TestHistory(t *testing.T) {
	h := NewHistory(3)
	posts := []*reddit.Post{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
		{ID: "4"},
		{ID: "5"},
	}

	if last := h.Contains(posts[0]); last != 0 {
		t.Error("Should be 0, got", last)
	}

	h.Push(posts[0])

	if last := h.Contains(posts[0]); last != 1 {
		t.Error("Should be 1, got", last)
	}

	h.Push(posts[1])
	if last := h.Contains(posts[0]); last != 2 {
		t.Error("Should be 2, got", last)
	}
	if last := h.Contains(posts[1]); last != 1 {
		t.Error("Should be 1, got", last)
	}

	h.Push(posts[2])
	if last := h.Contains(posts[0]); last != 3 {
		t.Error("Should be 2, got", last)
	}
	if last := h.Contains(posts[1]); last != 2 {
		t.Error("Should be 2, got", last)
	}
	if last := h.Contains(posts[2]); last != 1 {
		t.Error("Should be 1, got", last)
	}

	h.Push(posts[3])
	if last := h.Contains(posts[0]); last != 0 {
		t.Error("Should be 0, got", last)
	}
	if last := h.Contains(posts[1]); last != 3 {
		t.Error("Should be 3, got", last)
	}
	if last := h.Contains(posts[2]); last != 2 {
		t.Error("Should be 2, got", last)
	}
	if last := h.Contains(posts[3]); last != 1 {
		t.Error("Should be 1, got", last)
	}

	h.Push(posts[4])
	if last := h.Contains(posts[0]); last != 0 {
		t.Error("Should be 0, got", last)
	}
	if last := h.Contains(posts[1]); last != 0 {
		t.Error("Should be 0, got", last)
	}
	if last := h.Contains(posts[2]); last != 3 {
		t.Error("Should be 3, got", last)
	}
	if last := h.Contains(posts[3]); last != 2 {
		t.Error("Should be 2, got", last)
	}
	if last := h.Contains(posts[4]); last != 1 {
		t.Error("Should be 1, got", last)
	}
}
