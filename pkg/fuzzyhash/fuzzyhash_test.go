package fuzzyhash

import (
	"fmt"
	"testing"
)

func TestDifferences(t *testing.T) {
	for _, tc := range []struct {
		a    string
		b    string
		want int
	}{
		{
			a:    "cow",
			b:    "bow",
			want: 1,
		},
		{
			a:    "moo",
			b:    "mom",
			want: 1,
		},
		{
			a:    "abc",
			b:    "xyz",
			want: 3,
		},
		{
			a:    "aaa",
			b:    "aab",
			want: 1,
		},
		{
			a:    "aaaa",
			b:    "aab",
			want: 2,
		},
		{
			a:    "aaaa",
			b:    "aaa",
			want: 1,
		},
		{
			a:    "aaab",
			b:    "aaa",
			want: 1,
		},
	} {
		t.Run(fmt.Sprintf("%v_%v", tc.a, tc.b), func(t *testing.T) {
			a := New([]byte(tc.a))
			b := New([]byte(tc.b))
			diff := a.CompareTo(&b)
			if diff != tc.want {
				t.Errorf("CompareTo() = %v, want %v", diff, tc.want)
			}

			if a.CompareTo(&b) != b.CompareTo(&a) {
				t.Errorf("a.CompareTo(b) = %v, b.CompareTo(a) = %v", a.CompareTo(&b), b.CompareTo(&a))
			}
		})
	}
}
