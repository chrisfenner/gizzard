// Package hashfinder provides a class for searching hashes.
package hashfinder

import (
	"fmt"
	"hash"
	"sync"

	"github.com/chrisfenner/gizzard/pkg/fuzzyhash"
)

const maxIterations = 1000

// A HashFinder contains a list of hashes to search for, and can search all of
// them at once.
type HashFinder struct {
	// Tolerance is how many byte values can be different between two hashes.
	tolerance int
	// hashAlgs is all the hash algorithms the finder is checking.
	hashAlgs []func() hash.Hash
	// searchList is all the hashes this finder is searching for.
	searchList []fuzzyhash.FuzzyHash
}

// New creates a new HashFinder with the given tolerance, supporting the given
// algorithms.
func New(tolerance int, hashAlgs ...func() hash.Hash) HashFinder {
	return HashFinder{
		tolerance: tolerance,
		hashAlgs:  hashAlgs,
	}
}

// AddTarget adds a target hash to the HashFinder.
func (hf *HashFinder) AddTarget(data []byte) {
	hf.searchList = append(hf.searchList, fuzzyhash.New(data))
}

type worker = func(wg *sync.WaitGroup, data []byte, searchList []fuzzyhash.FuzzyHash, tolerance int, results chan string)

func makeHashRepeatedDataWorker(hashAlg func() hash.Hash) worker {
	return func(wg *sync.WaitGroup, data []byte, searchList []fuzzyhash.FuzzyHash, tolerance int, results chan string) {
		defer wg.Done()
		h := hashAlg()
		// Try hashing the data over and over.
		for i := 0; i < maxIterations; i++ {
			h.Write(data)
			candidate := fuzzyhash.New(h.Sum(nil))
			for _, target := range searchList {
				if candidate.CompareTo(&target) <= tolerance {
					results <- fmt.Sprintf("hash '%x' %d times", data, i)
					return
				}
			}
		}
	}
}

func makeRepeatedHashWorker(hashAlg func() hash.Hash) worker {
	return func(wg *sync.WaitGroup, data []byte, searchList []fuzzyhash.FuzzyHash, tolerance int, results chan string) {
		defer wg.Done()
		h := hashAlg()
		// Try hashing the hash of the data over and over.
		h.Write(data)
		data = h.Sum(nil)
		for i := 0; i < maxIterations; i++ {
			candidate := fuzzyhash.New(data)
			for _, target := range searchList {
				if candidate.CompareTo(&target) <= tolerance {
					results <- fmt.Sprintf("take the %d-th hash of '%x'", i, data)
					return
				}
			}
			h.Reset()
			h.Write(data)
			data = h.Sum(nil)
		}
	}
}

// Check checks to see if the given data matches any of the hashes this
// HashFinder is looking for.
func (hf *HashFinder) Check(data []byte) string {
	var workers []worker
	for _, alg := range hf.hashAlgs {
		workers = append(workers, makeHashRepeatedDataWorker(alg))
		workers = append(workers, makeRepeatedHashWorker(alg))
	}
	results := make(chan string, len(workers))
	var wg sync.WaitGroup
	wg.Add(len(workers))
	for _, worker := range workers {
		go worker(&wg, data, hf.searchList, hf.tolerance, results)
	}
	wg.Wait()
	if len(results) > 0 {
		return <-results
	}
	return ""
}
