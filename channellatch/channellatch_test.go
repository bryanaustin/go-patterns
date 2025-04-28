package channellatch

import (
	"context"
	"testing"
	"time"
)

const (
	TestChannelLatchDur = time.Millisecond * 50
)

// TestChannelLatchEmpty tests an empty ChannelLatch
func TestChannelLatchEmpty(t *testing.T) {
	t.Parallel()
	ctx, canc := context.WithTimeout(context.Background(), TestChannelLatchDur)
	defer canc()
	cl := New[int]()
	defer cl.Stop()
	c := cl.ChanDrainClose(ctx)
	helperChannelLatchNotReady(t, c)
	cl.Release()
	helperChannelLatchClosed(t, c)
}

// TestChannelLatchSingle tests a ChannelLatch with a single value
func TestChannelLatchSingle(t *testing.T) {
	t.Parallel()
	ctx, canc := context.WithTimeout(context.Background(), TestChannelLatchDur)
	defer canc()
	cl := New[int]()
	defer cl.Stop()
	c := cl.ChanDrainClose(ctx)
	cl.Add(1)
	helperChannelLatchNotReady(t, c)
	cl.Release()
	helperChannelLatchDrainLast(t, c, 1)
}

// TestChannelLatchTriple tests a ChannelLatch with three elements
func TestChannelLatchTriple(t *testing.T) {
	t.Parallel()
	ctx, canc := context.WithTimeout(context.Background(), TestChannelLatchDur)
	defer canc()
	cl := New[int]()
	defer cl.Stop()
	c := cl.ChanDrainClose(ctx)
	cl.Add(1)
	cl.Add(2)
	cl.Add(3)
	helperChannelLatchNotReady(t, c)
	cl.Release()
	helperChannelLatchDrainSingle(t, c, 1)
	helperChannelLatchDrainSingle(t, c, 2)
	helperChannelLatchDrainLast(t, c, 3)
}

// TestChannelLatchHold tests the ChannelLatch.Hold function
func TestChannelLatchHold(t *testing.T) {
	t.Parallel()
	cl := New[int]()
	defer cl.Stop()
	c := cl.Chan()
	cl.Add(1)
	cl.Add(2)
	cl.Add(3)
	helperChannelLatchNotReady(t, c)
	cl.Release()
	helperChannelLatchDrainSingle(t, c, 1)
	cl.Hold()
	helperChannelLatchNotReady(t, c)
	cl.Release()
	helperChannelLatchDrainSingle(t, c, 2)
	cl.Hold()
	helperChannelLatchNotReady(t, c)
	cl.Release()
	helperChannelLatchDrainSingle(t, c, 3)
}

// TestChannelLatchAddAfterDrain tests adding to the ChannelLatch after drain has started
func TestChannelLatchAddAfterDrain(t *testing.T) {
	t.Parallel()
	cl := New[int]()
	defer cl.Stop()
	ctx, canc := context.WithTimeout(context.Background(), TestChannelLatchDur)
	defer canc()
	c := cl.ChanDrainClose(ctx)
	cl.Add(1)
	cl.Add(2)
	cl.Release()
	cl.Hold() // using the "drainclose" function this will buffer one message
	helperChannelLatchDrainSingle(t, c, 1)
	helperChannelLatchNotReady(t, c)
	cl.Release()
	helperChannelLatchDrainLast(t, c, 2)
}

// TestChannelLatchRemove tests removing elements form the ChennelLatch
func TestChannelLatchRemove(t *testing.T) {
	t.Parallel()
	cl := New[int]()
	defer cl.Stop()
	cl.Release()
	ctx, canc := context.WithTimeout(context.Background(), TestChannelLatchDur)
	defer canc()
	c := cl.ChanDrainClose(ctx)
	cl.Add(1)
	cl.Add(2)
	cl.Add(3)
	cl.Remove(2)
	helperChannelLatchDrainSingle(t, c, 1)
	helperChannelLatchDrainLast(t, c, 3)
}

// TestChannelLatchRemoveNext tests removing the top most element from a ChannelLatch
func TestChannelLatchRemoveNext(t *testing.T) {
	t.Parallel()
	cl := New[int]()
	defer cl.Stop()
	cl.Release()
	c := cl.Chan()
	cl.Add(1)
	cl.Add(2)
	cl.Add(3)
	helperChannelLatchDrainSingle(t, c, 1)
	cl.Remove(2)
	helperChannelLatchDrainSingle(t, c, 3)
}

// TestChannelLatchRemoveNone tests removing an element from a ChannelLatch that doesn't exist
func TestChannelLatchRemoveNone(t *testing.T) {
	t.Parallel()
	cl := New[int]()
	defer cl.Stop()
	cl.Release()
	ctx, canc := context.WithTimeout(context.Background(), TestChannelLatchDur)
	defer canc()
	c := cl.ChanDrainClose(ctx)
	cl.Add(1)
	cl.Remove(3)
	helperChannelLatchDrainLast(t, c, 1)
}

func helperChannelLatchDrainSingle(t *testing.T, c <-chan int, expected int) {
	t.Helper()
	select {
	case v := <-c:
		if v != expected {
			t.Errorf("expected %d, got %d", expected, v)
		}
	}
}

func helperChannelLatchDrainLast(t *testing.T, c <-chan int, expected int) {
	t.Helper()
	select {
	case v := <-c:
		if v != expected {
			t.Errorf("expected %d, got %d", expected, v)
		}
	}
	helperChannelLatchClosed(t, c)
}

func helperChannelLatchNotReady(t *testing.T, c <-chan int) {
	t.Helper()
	select {
	case v := <-c:
		t.Errorf("channel should not have returned anything, got: %d", v)
	default:
	}
}

func helperChannelLatchClosed(t *testing.T, c <-chan int) {
	t.Helper()
	select {
	case _, ok := <-c:
		if ok {
			t.Error("channel should not be ok, but it was")
		}
	}
}
