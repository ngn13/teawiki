package util

import (
	"bytes"
)

type Buffer struct {
	Bytes []byte
	Index int
}

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

func (buff *Buffer) String() string {
	return string(buff.Bytes)
}

func (buff *Buffer) Length() int {
	return len(buff.Bytes)
}

func (buff *Buffer) Size() int {
	return buff.Index
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		Index: 0,
		Bytes: bytes.Repeat([]byte{0}, size),
	}
}
