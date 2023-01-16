package main

import (
	"fmt"
	"os"
	"testing"
)

func TestPrinter(t *testing.T) {
	p := &printer{w: os.Stdout}
	defer p.Close()

	fmt.Fprintln(p, "ciao")
	fmt.Fprintln(p, "ciao")
	fmt.Fprintln(p, "test")
	p.Println("hey")
	fmt.Fprintln(p, "hey")
	p.Println("hey")
	p.Println("hey")
	p.Printf("%s\n", "hey")
	fmt.Fprintln(p, "salut")
}
