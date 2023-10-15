package utils

import (
	"bytes"
	"crypto/rand"
	"github.com/sagernet/sing/common/buf"
	"io"
)

// GetRandomBytes is a substitute for the missing random.Default function on random.
func GetRandomBytes() (io.Reader, error) {
	randomData := make([]byte, 1024) // Change the size as needed
	_, err := rand.Read(randomData)
	if err != nil {
		return nil, err
	}
	return io.Reader(bytes.NewReader(randomData)), nil
}

// CopyBuffer  is a substitute for the missing Copy function on buf.Buffer.
func CopyBuffer(b *buf.Buffer) []byte {
	buffer := make([]byte, b.Len())
	copy(buffer, b.Bytes())
	return buffer
}
