package tsp

import (
	"bytes"
	_ "embed"
	"fmt"
	"testing"
)

//go:embed berlin52.tsp
var berlin52 []byte

func TestOpenFile(t *testing.T) {
	f, err := OpenFile(bytes.NewReader(berlin52))
	if err != nil {
		t.Fatal(err)
	}
	_ = f

	fmt.Printf("%+v", f)
}
