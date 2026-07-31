package notifications

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFanoutDeliversToSubscribers(t *testing.T) {
	ctx := context.Background()
	fanout := NewFanout(1)
	defer fanout.Close()

	first, unsubscribeFirst, err := fanout.Subscribe(ctx, 1)
	require.NoError(t, err)
	defer unsubscribeFirst()
	second, unsubscribeSecond, err := fanout.Subscribe(ctx, 1)
	require.NoError(t, err)
	defer unsubscribeSecond()

	notification := Notification{ID: "n1", UserID: "u1", Kind: "task.created", Payload: "task-1"}
	require.NoError(t, fanout.Publish(ctx, notification))

	assert.Equal(t, notification, receiveNotification(t, first))
	assert.Equal(t, notification, receiveNotification(t, second))
}

func TestFanoutCloseStopsCleanly(t *testing.T) {
	ctx := context.Background()
	fanout := NewFanout(0)
	ch, unsubscribe, err := fanout.Subscribe(ctx, 0)
	require.NoError(t, err)
	defer unsubscribe()

	done := make(chan struct{})
	go func() {
		fanout.Close()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("fanout close did not return")
	}

	_, ok := <-ch
	assert.False(t, ok)
	assert.ErrorIs(t, fanout.Publish(ctx, Notification{ID: "n2"}), ErrClosed)
}

func receiveNotification(t *testing.T, ch <-chan Notification) Notification {
	t.Helper()
	select {
	case notification := <-ch:
		return notification
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for notification")
	}
	return Notification{}
}
