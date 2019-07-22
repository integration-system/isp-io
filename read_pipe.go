package io

import "io"

type ReadPipe interface {
	io.ReadCloser
	Unshift(wr io.Reader)
	Last() io.Reader
}

type readChain struct {
	writers []io.Reader
}

func (p readChain) Close() error {
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

func (p readChain) Read(bytes []byte) (int, error) {
	if len(p.writers) == 0 {
		return 0, io.EOF
	}
	return p.writers[0].Read(bytes)
}

func (p *readChain) Unshift(wr io.Reader) {
	p.writers = append([]io.Reader{wr}, p.writers...)
}

func (p readChain) Last() io.Reader {
	if len(p.writers) == 0 {
		return nil
	}

	return p.writers[len(p.writers)-1]
}

func NewReadPipe(writers ...io.Reader) ReadPipe {
	return &readChain{
		writers: writers,
	}
}
