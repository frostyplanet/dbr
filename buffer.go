package dbr

import (
	"bytes"
	"sync"
)

type Buffer interface {
	WriteString(s string) (n int, err error)
	String() string

	WriteValue(v ...interface{}) (err error)
	Value() []interface{}
}

var bufferPool sync.Pool
var bufferVPool sync.Pool

type buffer struct {
	bytes.Buffer
	v []interface{}
}

func (b *buffer) Free() {
	if b.v != nil {
		bufferVPool.Put(b.v)
	}
	b.v = nil
	bufferPool.Put(b)
}

func NewBuffer() *buffer {
	b, ok := bufferPool.Get().(*buffer)
	if ! ok {
		b = new(buffer)
	} else {
		b.Reset()
	}
	b.v, ok = bufferVPool.Get().([]interface{})
	if b.v == nil {
		b.v = make([]interface{}, 0, 10)
	} else {
		b.v = b.v[0:0]
	}
	return b
}

func (b *buffer) WriteValue(v ...interface{}) error {
	b.v = append(b.v, v...)
	return nil
}

func (b *buffer) Value() []interface{} {
	return b.v
}
