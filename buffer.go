package main

import (
	"bytes"
	"fmt"
)

type Buffer struct {
	bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

func (b *Buffer) WriteF(format string, a ...any) {
	b.WriteString(fmt.Sprintf(format, a...))
}
func (b *Buffer) WriteS(format string) {
	b.WriteString(format)
}
