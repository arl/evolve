package bitstring

import "fmt"

func ExampleNew() {
	// create a 8 bits Bitstring
	bitstring, _ := New(8)
	// upon creation all bits are unset
	fmt.Println(bitstring)
	// Output: 00000000
}

func ExampleMakeFromString() {
	// create a Bitstring from string
	bitstring, _ := MakeFromString("101001")
	fmt.Println(bitstring)
	// Output: 101001
}

func ExampleBitstring_Len() {
	// create a 8 bits Bitstring
	bitstring, _ := New(8)
	fmt.Println(bitstring.Len())
	// Output: 8
}

func ExampleBitstring_Bit() {
	// create a 8 bits Bitstring
	bitstring, _ := New(8)
	fmt.Println(bitstring.Bit(7))
	// Output: false
}

func ExampleBitstring_SetBit() {
	// create a 8 bits Bitstring
	bitstring, _ := New(8)
	bitstring.SetBit(2, true)
	fmt.Println(bitstring)
	bitstring.SetBit(2, false)
	fmt.Println(bitstring)
	// Output: 00000100
	// 00000000
}

func ExampleBitstring_FlipBit() {
	// create a 8 bits Bitstring
	bitstring, _ := New(8)
	bitstring.FlipBit(2)
	fmt.Println(bitstring)
	// Output: 00000100
}

func ExampleBitstring_ZeroesCount() {
	// create a 8 bits Bitstring
	bitstring, _ := New(8)
	// upon creation all bits are unset
	fmt.Println(bitstring.ZeroesCount())
	// Output: 8
}

func ExampleBitstring_OnesCount() {
	// create a 8 bits Bitstring
	bitstring, _ := New(8)
	// upon creation all bits are unset
	fmt.Println(bitstring.OnesCount())
	// Output: 0
}

func ExampleBitstring_BigInt() {
	// create a 8 bits Bitstring
	bitstring, _ := MakeFromString("100")
	bi := bitstring.BigInt()
	fmt.Println(bi.Int64())
	// Output: 4
}

func ExampleSwapRange() {
	bs1, _ := MakeFromString("111")
	bs2, _ := MakeFromString("000")
	// starting from bit 2 of bs1, swap 1 bit with bs2
	SwapRange(bs1, bs2, 2, 1)
	fmt.Println(bs1)
	// Output: 011
}
