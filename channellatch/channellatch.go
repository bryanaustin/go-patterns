package channellatch

import (
	"context"
)

// ChannelLatch allows you to queue elements to be processed on release
type ChannelLatch[T comparable] struct {
	data        []T
	chanAdd     chan T
	chanRm      chan T
	chanDrain   chan T
	chanDrained chan struct{}
	chanRelease chan struct{}
	chanHold    chan struct{}
	chanClose   chan struct{}
}

func New[T comparable]() *ChannelLatch[T] {
	result := &ChannelLatch[T]{
		chanAdd:     make(chan T),
		chanRm:      make(chan T),
		chanDrain:   make(chan T),
		chanDrained: make(chan struct{}),
		chanRelease: make(chan struct{}),
		chanHold:    make(chan struct{}),
		chanClose:   make(chan struct{}),
	}
	go result.run()
	return result
}

func (cl *ChannelLatch[T]) run() {
	defer close(cl.chanDrain)
	defer close(cl.chanDrained)

	var (
		next         T // next value queued to be released
		empty        bool
		skipHold     bool
		skipReleased bool
	)

	for {

		// Holding
		if !skipHold {
			for holding := true; holding; {
				select {
				case ad := <-cl.chanAdd:
					cl.data = append([]T{ad}, cl.data...)
				case rm := <-cl.chanRm:
					cl.doRm(rm)
				case <-cl.chanRelease:
					holding = false
				case <-cl.chanHold:
					// do nothing, already holding
				case <-cl.chanClose:
					return
				}
			}
		}

		skipHold = false
		next, empty = cl.pop()

		// Releasing
		for releasing := true; releasing && !empty; {
			select {
			case ad := <-cl.chanAdd:
				cl.data = append([]T{ad}, cl.data...)
			case rm := <-cl.chanRm:
				if rm == next {
					next, empty = cl.pop()
				}
				cl.doRm(rm)
			case cl.chanDrain <- next:
				next, empty = cl.pop()
			case <-cl.chanHold:
				cl.data = append(cl.data, next)
				releasing = false
				skipReleased = true
			case <-cl.chanRelease:
				// do nothing, already releasing
			case <-cl.chanClose:
				return
			}
		}

		// All released
		if !skipReleased {
			for released := true; released; {
				select {
				case ad := <-cl.chanAdd:
					cl.data = append([]T{ad}, cl.data...)
					released = false
					skipHold = true
				case <-cl.chanRm:
					// nothing to remove
				case cl.chanDrained <- struct{}{}:
				case <-cl.chanHold:
					released = false
				case <-cl.chanRelease:
					// do nothing, already released
				case <-cl.chanClose:
					return
				}
			}
		}
		skipReleased = false
	}
}

func (cl *ChannelLatch[T]) pop() (vaule T, empty bool) {
	empty = len(cl.data) < 1
	if !empty {
		vaule = cl.data[len(cl.data)-1]
		cl.data = cl.data[:len(cl.data)-1]
	}
	return
}

// doRm removes the provided element from stored data
func (cl *ChannelLatch[T]) doRm(x T) (found bool) {
	var dummy T
	for i := 0; i < len(cl.data); {
		if cl.data[i] == x {
			// remove
			if i < len(cl.data)-1 {
				copy(cl.data[i:], cl.data[i+1:])
			}
			cl.data[len(cl.data)-1] = dummy // release for GC if pointer
			cl.data = cl.data[:len(cl.data)-1]
		} else {
			i++
		}
	}
	return
}

// Add will have ChannelLatch hold on to this element until released.
// Thread Safe.
func (cl *ChannelLatch[T]) Add(x T) {
	cl.chanAdd <- x
}

// Remove will have ChannelLatch remove all matching elements that it currently has held.
// Thread Safe.
func (cl *ChannelLatch[T]) Remove(x T) {
	cl.chanRm <- x
}

// Release will allow processing of data.
// Thread Safe.
func (cl *ChannelLatch[T]) Release() {
	cl.chanRelease <- struct{}{}
}

// Hold will stop the processing of data.
// Thread Safe.
func (cl *ChannelLatch[T]) Hold() {
	cl.chanHold <- struct{}{}
}

// WaitDrined will wait until ChannelLatch hits the drained or closed states
func (cl *ChannelLatch[T]) WaitDrined() {
	<-cl.chanDrained
}

// Stop will end all goroutines.
// Thread Safe.
func (cl *ChannelLatch[T]) Stop() {
	close(cl.chanClose)
}

// Chan returns a channel that will get data when released.
// Thread Safe.
func (cl *ChannelLatch[T]) Chan() <-chan T {
	return cl.chanDrain
}

// ChanDrainClose returns a channel that will drain data when released and close after.
// This also buffers one item at a time, it will still return one value after you've run Hold.
// Thread Safe.
func (cl *ChannelLatch[T]) ChanDrainClose(ctx context.Context) <-chan T {
	result := make(chan T)
	go func() {
		defer close(result)
		for {
			select {
			case x := <-cl.chanDrain:
				result <- x
			case <-ctx.Done():
				return
			case <-cl.chanDrained:
				return
			}
		}
	}()
	return result
}
