package bitstring

import "fmt"

func ExampleNew() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	// upon creation all bits are unset
	fmt.Println(bitstring)
	// Output: 00000000
}

func ExampleMakeFromString() {
	// create a BitString from string
	bitstring, _ := MakeFromString("101001")
	fmt.Println(bitstring)
	// Output: 101001
}

func ExampleBitString_Len() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	fmt.Println(bitstring.Len())
	// Output: 8
}

func ExampleBitString_Bit() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	fmt.Println(bitstring.Bit(7))
	// Output: false
}

func ExampleBitString_SetBit() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	bitstring.SetBit(2, true)
	fmt.Println(bitstring)
	bitstring.SetBit(2, false)
	fmt.Println(bitstring)
	// Output: 00000100
	// 00000000
}

func ExampleBitString_FlipBit() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	bitstring.FlipBit(2)
	fmt.Println(bitstring)
	// Output: 00000100
}

func ExampleBitString_ZeroesCount() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	// upon creation all bits are unset
	fmt.Println(bitstring.ZeroesCount())
	// Output: 8
}

func ExampleBitString_OnesCount() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	// upon creation all bits are unset
	fmt.Println(bitstring.OnesCount())
	// Output: 0
}

func ExampleBitString_BigInt() {
	// create a 8 bits BitString
	bitstring, _ := MakeFromString("100")
	bi := bitstring.BigInt()
	fmt.Println(bi.Int64())
	// Output: 4
}

func ExampleBitString_SwapRange() {
	bs1, _ := MakeFromString("111")
	bs2, _ := MakeFromString("000")
	// starting from bit 2 of bs1, swap 1 bit with bs2
	bs1.SwapRange(bs2, 2, 1)
	fmt.Println(bs1)
	// Output: 011
}
