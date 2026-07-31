package notifications

import (
	"context"
	"errors"
	"sync"
)

// ErrClosed is returned when publishing or subscribing after shutdown.
var ErrClosed = errors.New("notifications fanout closed")

// Fanout broadcasts notifications to subscribed channels from one dispatcher goroutine.
type Fanout struct {
	publish    chan Notification
	register   chan chan Notification
	unregister chan chan Notification
	stop       chan struct{}
	closed     chan struct{}
	once       sync.Once
}

// NewFanout starts a dispatcher with the given channel buffer size.
func NewFanout(buffer int) *Fanout {
	if buffer < 0 {
		buffer = 0
	}
	f := &Fanout{
		publish:    make(chan Notification, buffer),
		register:   make(chan chan Notification),
		unregister: make(chan chan Notification),
		stop:       make(chan struct{}),
		closed:     make(chan struct{}),
	}
	go f.run()
	return f
}

// Subscribe registers a notification channel and returns an unsubscribe function.
func (f *Fanout) Subscribe(ctx context.Context, buffer int) (<-chan Notification, func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	if buffer < 0 {
		buffer = 0
	}
	ch := make(chan Notification, buffer)
	select {
	case f.register <- ch:
		var once sync.Once
		unsubscribe := func() {
			once.Do(func() {
				select {
				case f.unregister <- ch:
				case <-f.closed:
				}
			})
		}
		return ch, unsubscribe, nil
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	case <-f.closed:
		return nil, nil, ErrClosed
	}
}

// Publish queues a notification for asynchronous delivery.
func (f *Fanout) Publish(ctx context.Context, notification Notification) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case f.publish <- notification:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-f.closed:
		return ErrClosed
	}
}

// Close stops the dispatcher and closes all subscriber channels.
func (f *Fanout) Close() {
	f.once.Do(func() {
		close(f.stop)
		<-f.closed
	})
}

func (f *Fanout) run() {
	subscribers := make(map[chan Notification]struct{})
	defer func() {
		for ch := range subscribers {
			close(ch)
		}
		close(f.closed)
	}()

	for {
		select {
		case ch := <-f.register:
			subscribers[ch] = struct{}{}
		case ch := <-f.unregister:
			if _, ok := subscribers[ch]; ok {
				delete(subscribers, ch)
				close(ch)
			}
		case notification := <-f.publish:
			for ch := range subscribers {
				select {
				case ch <- notification:
				default:
				}
			}
		case <-f.stop:
			return
		}
	}
}
