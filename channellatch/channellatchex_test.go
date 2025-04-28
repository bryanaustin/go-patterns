package channellatch

import (
	"context"
	"fmt"
	"sync"
)

func Example() {
	var wg sync.WaitGroup
	cl := New[int]()
	defer cl.Stop()

	go func() {
		defer wg.Done()
		for x := range cl.ChanDrainClose(context.Background()) {
			fmt.Printf(" %d", x)
		}
	}()

	// Add and remove
	cl.Add(1)
	cl.Add(2)
	cl.Add(3)
	cl.Add(4)
	cl.Remove(2)
	cl.Remove(1)

	// Prepare for go routine process
	wg.Add(1)

	// Prove that the go routine hasn't printed anything yet
	fmt.Print("9")

	// Release the ChannelLatch
	cl.Release()
	wg.Wait()

	// Output:
	// 9 3 4
}
