package utils

import "sync"

type Broadcaster[T any] struct {
	subs []chan T
	mu   sync.Mutex
}

func NewBroadcaster[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{
		subs: make([]chan T, 0),
	}
}

func (b *Broadcaster[T]) Subscribe() <-chan T {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan T)
	b.subs = append(b.subs, ch)

	return ch
}

func (b *Broadcaster[T]) Broadcast(val T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, ch := range b.subs {
		select {
		case ch <- val:
		default:
			continue
		}
	}
}
