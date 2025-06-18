package util

import (
	"bytes"
	"fmt"
	"io"
)

// buffer is used to store a specific amount of bytes, if the buffer is full,
// adding more data to it will trim the start bytes
type Buffer struct {
	Bytes []byte
	Index int
}

// clears the contents of the buffer and resets it's position
func (buff *Buffer) Clear() {
	for i := range buff.Index {
		buff.Bytes[i] = 0
	}

	buff.Index = 0
}

// pushes a single byte to the buffer
func (buff *Buffer) Push(b byte) int {
	if buff.Index >= len(buff.Bytes) {
		copy(buff.Bytes, buff.Bytes[1:])
		buff.Bytes[buff.Index-1] = b
	} else {
		buff.Bytes[buff.Index] = b
		buff.Index++
	}

	return buff.Index
}

// converts buffer to string
func (buff *Buffer) String() string {
	return string(buff.Bytes[:buff.Index])
}

// get the total size of the buffer
func (buff *Buffer) Size() int {
	return len(buff.Bytes)
}

// returns the current length of the buffer
func (buff *Buffer) Len() int {
	return buff.Index
}

// pushes total of n bytes from the specified reader into the buffer
func (buff *Buffer) From(r io.Reader, n int) error {
	if n > len(buff.Bytes) {
		return fmt.Errorf("size is too large for buffer")
	}

	if buff.Index >= len(buff.Bytes) {
		copy(buff.Bytes, buff.Bytes[n:])
		buff.Index -= 1
	}

	reader := io.LimitReader(r, int64(n))
	read, err := reader.Read(buff.Bytes[buff.Index:])

	if err != nil {
		return err
	}

	if read != n {
		return fmt.Errorf("failed to read all %d bytes", n)
	}

	buff.Index += read
	return nil
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		Index: 0,
		Bytes: bytes.Repeat([]byte{0}, size),
	}
}
