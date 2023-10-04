package chunkbuffer

import (
	"io"
)

type ChunkBuffer struct {
	chunks                         [][]byte
	readindex, writeindex, readsub int
	length                         int
}

// Read implements io.Reader
func (b *ChunkBuffer) Read(ot []byte) (n int, err error) {
	on := len(ot)
	for {
		if b.readindex == b.writeindex {
			b.readindex = 0
			b.writeindex = 0
			b.readsub = 0
			err = io.EOF
			break // Empty
		}
		if n == on {
			break // Full
		}

		w := copy(ot[n:], b.chunks[b.readindex][b.readsub:])
		n += w
		b.readsub += w

		if b.readsub == len(b.chunks[b.readindex]) {
			b.chunks[b.readindex] = nil
			b.readindex = (b.readindex + 1) % len(b.chunks)
			b.readsub = 0
			if b.readindex == b.writeindex && b.readindex > 0 {
				b.readindex = 0
				b.writeindex = 0
			}
		}
		if w < 1 {
			break
		}
	}
	b.length -= n
	return
}

// Write implements io.Writer
func (b *ChunkBuffer) Write(in []byte) (n int, err error) {
	nu := append([]byte(nil), in...) // copy
	n = len(nu)
	b.length += n

	// Inintal alloc and ideal allocs
	if b.writeindex == len(b.chunks) {
		// Expand and copy
		nuchunks := make([][]byte, max(len(b.chunks)*2,8))
		copy(nuchunks, b.chunks)
		b.chunks = nuchunks
	}

	b.chunks[b.writeindex] = nu
	b.writeindex = (b.writeindex + 1) % len(b.chunks)

	if b.writeindex == b.readindex {
		if b.readindex > 0 {
			// Just filled, rewrite chunks with more space
			nuchunks := make([][]byte, max(len(b.chunks)*2,8))
			f := copy(nuchunks[0:], b.chunks[b.readindex:])
			a := copy(nuchunks[f:], b.chunks[:b.writeindex])
			b.readindex = 0
			b.writeindex = f + a
			b.chunks = nuchunks
		} else {
			// Rewrite not needed, just prep for ideal alloc next round
			b.writeindex = len(b.chunks)
		}
	}
	return
}

// Len is the number of bytes held in the buffer
func (b ChunkBuffer) Len() int {
	return b.length
}
