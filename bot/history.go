package bot

import "github.com/turnage/graw/reddit"

type History struct {
	length int
	ids    []string
	prev   int
	next   int
}

func NewHistory(length int) *History {
	return &History{
		length: length,
		prev:   length,
		next:   length,
		ids:    make([]string, length),
	}
}

func (h *History) Push(p *reddit.Post) {
	h.ids[h.next%h.length] = p.ID
	h.prev = h.next
	h.next += 1
}

func (h *History) Contains(p *reddit.Post) int {
	end := h.prev - h.length
	for i := h.prev; i > end; i-- {
		if h.ids[i%h.length] == p.ID {
			return int(h.next - i)
		}
	}

	return 0
}
