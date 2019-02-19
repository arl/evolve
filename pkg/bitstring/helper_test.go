package bitstring

import "fmt"

// returns a string representing the first n bits of the base-2 representation
// of x (unsigned).
func sprintubits(val uint32, nbits uint) string {
	return fmt.Sprintf(fmt.Sprintf("%%0%db", nbits), val)
}

// returns a string representing the first n bits of the base-2 representation
// of x (signed).
func sprintsbits(val int32, nbits uint) string {
	if val < 0 {
		// casting to uint will show us the 2's complement
		return sprintubits(uint32(val), nbits)
	}
	return fmt.Sprintf(fmt.Sprintf("%%0%db", nbits), val)
}
