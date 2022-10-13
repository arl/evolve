package main

import "testing"

func Test_funcstack(t *testing.T) {
	got := ""
	var _defer funcstack
	_defer.add(func() { got += "c" })
	_defer.add(func() { got += "b" })
	_defer.add(func() { got += "a" })

	_defer.run()

	want := "abc"
	if want != got {
		t.Errorf("got %s, want %s", got, want)
	}
}
