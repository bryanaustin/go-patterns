package chunkbuffer

import (
	"testing"
	"bytes"
	"io"
	"github.com/google/go-cmp/cmp"
)

var numbs = make([]byte, 1024)

func init() {
	for i := 0; i < len(numbs); i++ {
		numbs[i] = byte(i)
	}
}

func TestBasic(t *testing.T) {
	t.Parallel()
	b := new(ChunkBuffer)
	b.Write(numbs[:64])
	out := make([]byte, 64)
	b.Read(out)
	if msg := cmp.Diff(numbs[:64], out); len(msg) > 0 {
		t.Error(msg)
	}
}

func TestHalfRead(t *testing.T) {
	t.Parallel()
	b := new(ChunkBuffer)
	b.Write(numbs[:64])
	verifyState(t, *b, 0, 0, 1, 64)
	out1 := make([]byte, 32)
	b.Read(out1)
	verifyState(t, *b, 0, 32, 1, 32)
	out2 := make([]byte, 32)
	b.Read(out2)
	verifyState(t, *b, 0, 0, 0, 0)
	if msg := cmp.Diff(numbs[:32], out1); len(msg) > 0 {
		t.Error(msg)
	}
	if msg := cmp.Diff(numbs[32:64], out2); len(msg) > 0 {
		t.Error(msg)
	}
}

func TestOverRead(t *testing.T) {
	t.Parallel()
	b := new(ChunkBuffer)
	b.Write(numbs[:64])
	verifyState(t, *b, 0, 0, 1, 64)
	out1 := make([]byte, 32)
	b.Read(out1)
	verifyState(t, *b, 0, 32, 1, 32)
	out2 := make([]byte, 64)
	n, _ := b.Read(out2)
	out2 = out2[:n]
	verifyState(t, *b, 0, 0, 0, 0)
	if msg := cmp.Diff(numbs[:32], out1); len(msg) > 0 {
		t.Error(msg)
	}
	if msg := cmp.Diff(numbs[32:64], out2); len(msg) > 0 {
		t.Error(msg)
	}
}

func BenchmarkBytesBufferTest1(b *testing.B) {
	f := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		test1(f, io.Discard)
	}
}

func BenchmarkChunkBufferTest1(b *testing.B) {
	f := new(ChunkBuffer)
	for i := 0; i < b.N; i++ {
		test1(f, io.Discard)
	}
}

func test1(x io.ReadWriter, w io.Writer) {
	for i := 2; i < len(numbs); i++ {
		x.Write(numbs[0:i])
		io.CopyN(w, x, int64(i-1))
	}
}

func BenchmarkBytesBufferTest2(b *testing.B) {
	f := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		test2(f, io.Discard)
	}
}

func BenchmarkChunkBufferTest2(b *testing.B) {
	f := new(ChunkBuffer)
	for i := 0; i < b.N; i++ {
		test2(f, io.Discard)
	}
}

func test2(x io.ReadWriter, w io.Writer) {
	var count int64
	for i := 1; i < len(numbs)*64; i++ {
		n, _ := x.Write(numbs[:i%len(numbs)])
		count += int64(n)
		cn := count - int64(i) - 1
		var m int64
		if cn > 0 {
			m, _ = io.CopyN(w, x, cn)
		}
		count -= m
	}
}
func BenchmarkBytesBufferTest3(b *testing.B) {
	f := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		test3(f, io.Discard)
	}
}

func BenchmarkChunkBufferTest3(b *testing.B) {
	f := new(ChunkBuffer)
	for i := 0; i < b.N; i++ {
		test3(f, io.Discard)
	}
}

func test3(x io.ReadWriter, w io.Writer) {
	for i := 1; i < len(numbs); i++ {
		x.Write(numbs[0:(i)])
		io.CopyN(w, x, int64(i))
	}
}

func verifyState(t *testing.T, b ChunkBuffer, ri, rs, wi, l int) {
	t.Helper()

	if msg := cmp.Diff(b.readindex, ri); len(msg) > 0 {
		t.Error(msg)
	}

	if msg := cmp.Diff(b.readsub, rs); len(msg) > 0 {
		t.Error(msg)
	}

	if msg := cmp.Diff(b.writeindex, wi); len(msg) > 0 {
		t.Error(msg)
	}

	if msg := cmp.Diff(b.length, l); len(msg) > 0 {
		t.Error(msg)
	}
}
