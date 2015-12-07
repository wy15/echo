package bpool

import "bytes"

type SizedBufferPool struct {
	c chan *bytes.Buffer
	a int
}

func NewSizedBufferPool(size, alloc int) *SizedBufferPool {
	return &SizedBufferPool{
		c: make(chan *bytes.Buffer, size),
		a: alloc,
	}
}

func (bp *SizedBufferPool) Get() (b *bytes.Buffer) {
	select {
	case b = <-bp.c:
	default:
		b = bytes.NewBuffer(make([]byte, 0, bp.a))
	}

	return
}

func (bp *SizedBufferPool) Put(b *bytes.Buffer) {
	b.Reset()

	if b.Cap() > bp.a {
		b = bytes.NewBuffer(make([]byte, 0, bp.a))
	}

	select {
	case bp.c <- b:
	default:
	}
}
