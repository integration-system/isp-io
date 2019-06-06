package io

import (
	"io"
)

type Pipe interface {
	io.WriteCloser
	Unshift(wr io.Writer)
	Last() io.Writer
}

type writeChain struct {
	writers []io.Writer
}

func (p writeChain) Close() error {
	for _, w := range p.writers {
		if f, ok := w.(interface{ Flush() error }); ok {
			f.Flush() //TODO handle error
		}
		if c, ok := w.(io.Closer); ok {
			c.Close() //TODO handle error
		}
	}
	return nil
}

func (p writeChain) Write(bytes []byte) (int, error) {
	if len(p.writers) == 0 {
		return 0, nil
	}
	return p.writers[0].Write(bytes)
}

func (p *writeChain) Unshift(wr io.Writer) {
	p.writers = append([]io.Writer{wr}, p.writers...)
}

func (p writeChain) Last() io.Writer {
	if len(p.writers) == 0 {
		return nil
	}

	return p.writers[len(p.writers)-1]
}

func NewPipe(writers ...io.Writer) Pipe {
	return &writeChain{
		writers: writers,
	}
}
