package io

import "io"

type ReadPipe interface {
	io.ReadCloser
	Unshift(wr io.Reader)
	Last() io.Reader
}

type readChain struct {
	readers []io.Reader
}

func (p readChain) Close() error {
	for _, w := range p.readers {
		if f, ok := w.(interface{ Flush() error }); ok {
			if err := f.Flush(); err != nil {
				return err
			}
		}
		if f, ok := w.(interface{ Flush() }); ok {
			f.Flush()
		}
		if c, ok := w.(io.Closer); ok {
			if err := c.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p readChain) Read(bytes []byte) (int, error) {
	if len(p.readers) == 0 {
		return 0, io.EOF
	}
	return p.readers[0].Read(bytes)
}

func (p *readChain) Unshift(wr io.Reader) {
	p.readers = append([]io.Reader{wr}, p.readers...)
}

func (p readChain) Last() io.Reader {
	if len(p.readers) == 0 {
		return nil
	}

	return p.readers[len(p.readers)-1]
}

func NewReadPipe(writers ...io.Reader) ReadPipe {
	return &readChain{
		readers: writers,
	}
}
