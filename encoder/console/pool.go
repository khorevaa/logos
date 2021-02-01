package console

import (
	"go.uber.org/zap/buffer"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} {
		return &Encoder{}
	},
}

func get() *Encoder {
	return pool.Get().(*Encoder)
}

func put(enc *Encoder) {
	enc.buf = nil
	pool.Put(enc)
}

var bufferpool = buffer.NewPool()
