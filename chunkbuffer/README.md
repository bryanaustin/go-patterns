# ChunkBuffer
This is an alternate to the `bytes.Buffer` type that implements `io.ReadWriter` and uses a slice of byte slices instead of the one single byte slice that `bytes.Buffer` uses. The hope was that it would have fewer copies and handle partial drains better; however, when benchmarking it, it was always slower and sometimes used less memory.

I may come back to this one day to better graph its performance compared to `bytes.Buffer`, but it seems too niche to try to convince you to use it right now.

Here is the last set of benchmarks I ran:
```
$ go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/bryanaustin/go-patterns/chunkbuffer
cpu: AMD Ryzen Threadripper 2920X 12-Core Processor 
BenchmarkBytesBufferTest1-24    	    4762	    564921 ns/op	  848905 B/op	    1022 allocs/op
BenchmarkChunkBufferTest1-24    	    2012	    711161 ns/op	  580746 B/op	    2044 allocs/op
BenchmarkBytesBufferTest2-24    	      58	  34713034 ns/op	50338368 B/op	   65380 allocs/op
BenchmarkChunkBufferTest2-24    	      31	  44978594 ns/op	37172105 B/op	  130874 allocs/op
BenchmarkBytesBufferTest3-24    	    4701	    257929 ns/op	   24619 B/op	    1023 allocs/op
BenchmarkChunkBufferTest3-24    	    2409	    478202 ns/op	  580080 B/op	    2046 allocs/op
```
