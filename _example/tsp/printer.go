package main

import (
	"fmt"
	"io"
)

// printer prints to w, counting bursts.
type printer struct {
	prev  string
	count int
	w     io.Writer
}

func (p *printer) Close() error {
	if p.count == 0 {
		fmt.Fprint(p.w, p.prev)
	} else {
		fmt.Fprintf(p.w, "[%d] %s\n", p.count+1, p.prev)
	}
	return nil
}

func (p *printer) Println(a ...any) (n int, err error) {
	return fmt.Fprintln(p, a...)
}

func (p *printer) Printf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(p, format, a...)
}

func (p *printer) Write(b []byte) (n int, err error) {
	switch {
	case string(b) == p.prev:
		p.count++
		return len(b), nil
	case p.count != 0:
		n, err := fmt.Fprintf(p.w, "[%d] %s", p.count+1, p.prev)
		p.prev = string(b)
		p.count = 0
		return n, err
	}
	fmt.Fprintf(p.w, p.prev)
	p.prev = string(b)
	p.count = 0
	return len(b), nil
}
