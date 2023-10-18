// Package bytefield contains a helper type for dealing with unordered collections of bytes.
package bytefield

// A Bytefield represents a hash as counts of values.
type Bytefield struct {
	// How many times each byte value has appeared.
	counts [256]int
}

// CountOf returns the number of times the given value appears in the Bytefield.
func (b *Bytefield) CountOf(value byte) int {
	return b.counts[value]
}

// Increment sets the bit for the given value.
func (b *Bytefield) Increment(value byte) {
	b.counts[value]++
}
