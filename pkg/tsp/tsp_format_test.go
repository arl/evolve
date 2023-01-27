package tsp

import (
	"bytes"
	_ "embed"
	"fmt"
	"testing"
)

//go:embed testdata/berlin52.tsp
var berlin52 []byte

func TestLoad_berlin52(t *testing.T) {
	f, err := Load(bytes.NewReader(berlin52))
	if err != nil {
		t.Fatal(err)
	}
	_ = f

	fmt.Printf("%+v", f)
}

// optimum 629
//
//go:embed testdata/eil101.tsp
var eil101 []byte

func TestLoad_eil101(t *testing.T) {
	f, err := Load(bytes.NewReader(eil101))
	if err != nil {
		t.Fatal(err)
	}
	_ = f

	fmt.Printf("%+v", f)
}
