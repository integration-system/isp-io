package io

import (
	"io"
)

type Pipe []io.Writer

func (p Pipe) Close() error {
	for _, w := range p {
		if f, ok := w.(interface{ Flush() error }); ok {
			f.Flush() //TODO handle error
		}
		if c, ok := w.(io.Closer); ok {
			c.Close() //TODO handle error
		}
	}
	return nil
}

func (p Pipe) Write(bytes []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	return p[0].Write(bytes)
}

func (p *Pipe) Unshift(wr io.Writer) {
	pipe := append(Pipe{wr}, *p...)
	p = &pipe
}

func (p Pipe) Last() io.Writer {
	if len(p) == 0 {
		return nil
	}

	return p[len(p)-1]
}

func NewPipe() Pipe {
	return Pipe{}
}
