package playback

import (
	"sync"
	"time"
)

type Buffer struct {
	s   Status
	mux sync.Mutex
}

func (b *Buffer) Push(s Status) {
	b.mux.Lock()
	defer b.mux.Unlock()
	if b.s == nil || time.Since(s.Created()) < time.Since(b.s.Created()) {
		b.s = s
	}
}

func (b *Buffer) Pop() Status {
	b.mux.Lock()
	defer b.mux.Unlock()
	s := b.s
	b.s = nil
	return s
}
func (b *Buffer) Peek() Status {
	return b.s
}
