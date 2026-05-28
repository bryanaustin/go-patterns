
import (
	"encoding/ascii85"
	"encoding/base64"
	"encoding/hex"
	"math/rand/v2"
)

// RandStringASCII85 returns a random ASCII85-encoded string of exactly length n.
func RandStringASCII85(n int) string {
	if n <= 0 {
		return ""
	}
	// 5 ascii85 chars encode 4 raw bytes; over-generate then trim.
	raw := make([]byte, (n*4+4)/5)
	for i := range raw {
		raw[i] = byte(rand.UintN(256))
	}
	dst := make([]byte, ascii85.MaxEncodedLen(len(raw)))
	m := ascii85.Encode(dst, raw)
	return string(dst[:m])[:n]
}

// RandStringBase64 returns a random base64 (url, padded) string of exactly length n.
func RandStringBase64(n int) string {
	if n <= 0 {
		return ""
	}
	// 4 base64 chars per 3 raw bytes; over-generate then trim.
	raw := make([]byte, (n*3+3)/4)
	for i := range raw {
		raw[i] = byte(rand.UintN(256))
	}
	return base64.URLEncoding.EncodeToString(raw)[:n]
}

// RandStringHex returns a random hex string of exactly length n.
func RandStringHex(n int) string {
	if n <= 0 {
		return ""
	}
	// 2 hex chars per byte.
	raw := make([]byte, (n+1)/2)
	for i := range raw {
		raw[i] = byte(rand.UintN(256))
	}
	return hex.EncodeToString(raw)[:n]
}
