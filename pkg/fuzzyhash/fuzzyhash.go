// Package fuzzyhash contains a helper type for comparing hash values fuzzily.
package fuzzyhash

import (
	"fmt"

	"github.com/chrisfenner/gizzard/pkg/bytefield"
)

// FuzzyHash encapsulates a hash value, with fuzzy-matching capabilities.
type FuzzyHash struct {
	value []byte
	field bytefield.Bytefield
}

// New creates a new FuzzyHash out of the given value.
func New(value []byte) FuzzyHash {
	var result FuzzyHash
	result.value = make([]byte, len(value))
	copy(result.value, value)
	for _, b := range result.value {
		result.field.Increment(b)
	}
	return result
}

// String implements the Stringer interface.
func (fh *FuzzyHash) String() string {
	return fmt.Sprintf("%x", fh.value)
}

// CompareTo returns how many byte values were different between the two hashes
// (ignoring order).
func (fh *FuzzyHash) CompareTo(other *FuzzyHash) int {
	// Total up the absolute values of the differences of the byte counts of
	// all the byte values.
	difference := 0
	for i := 0; i < 256; i++ {
		thisCount := fh.field.CountOf(byte(i))
		thatCount := other.field.CountOf(byte(i))

		if thisCount == thatCount {
			continue
		} else if thisCount < thatCount {
			difference += (thatCount - thisCount)
		} else {
			difference += (thisCount - thatCount)
		}
	}

	// Differences are always doubled when counting with the algorithm above.
	// i.e., if one byte is wrong, then one count will be too low and another
	// will be too high.
	difference /= 2

	// If the hashes have different sizes, add the difference in sizes.
	thisLen := len(fh.value)
	thatLen := len(other.value)
	if thisLen != thatLen {
		if thisLen < thatLen {
			difference += (thatLen - thisLen)
		} else {
			difference += (thisLen - thatLen)
		}
	}

	return difference
}
