package hashfinder

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"testing"

	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
)

// TODO this test is pretty quick and dirty. It could be both more succinct and comprehensive.
func TestHashFinder(t *testing.T) {
	for _, tc := range []struct {
		name string
		alg  func() hash.Hash
	}{
		{
			"SHA1",
			sha1.New,
		},
		{
			"SHA256",
			sha256.New,
		},
		{
			"MD5",
			md5.New,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("hash once", func(t *testing.T) {
				sample := tc.alg()
				sample.Write([]byte("test"))
				sampleHash := sample.Sum(nil)

				hf := New(0, tc.alg)
				hf.AddTarget(sampleHash)
				result := hf.Check([]byte("test"))
				if result == "" {
					t.Errorf("Check() = '', want success")
				}

				// mess up one byte
				sampleHash[0] ^= 0xff
				hf = New(0, tc.alg)
				hf.AddTarget(sampleHash)
				result = hf.Check([]byte("test"))
				if result != "" {
					t.Errorf("Check() = %q, want ''", result)
				}

				hf = New(1, tc.alg)
				hf.AddTarget(sampleHash)
				result = hf.Check([]byte("test"))
				if result == "" {
					t.Errorf("Check() = '', want success")
				}

				// mess up 4 bytes
				sampleHash[1] ^= 0xff
				sampleHash[2] ^= 0xff
				sampleHash[3] ^= 0xff
				hf = New(4, tc.alg)
				hf.AddTarget(sampleHash)
				result = hf.Check([]byte("test"))
				if result == "" {
					t.Errorf("Check() = '', want success")
				}
				hf = New(0, tc.alg)
				hf.AddTarget(sampleHash)
				result = hf.Check([]byte("test"))
				if result != "" {
					t.Errorf("Check() = %q, want ''", result)
				}
			})
			t.Run("hash squared", func(t *testing.T) {
				sample := tc.alg()
				sample.Write([]byte("test"))
				sampleHash := sample.Sum(nil)
				// hash it again
				sample.Reset()
				sample.Write(sampleHash)
				sampleHash = sample.Sum(nil)

				hf := New(0, tc.alg)
				hf.AddTarget(sampleHash)
				result := hf.Check([]byte("test"))
				if result == "" {
					t.Errorf("Check() = '', want success")
				}

				// mess up one byte
				sampleHash[0] ^= 0xff
				hf = New(0, tc.alg)
				hf.AddTarget(sampleHash)
				result = hf.Check([]byte("test"))
				if result != "" {
					t.Errorf("Check() = %q, want ''", result)
				}

				hf = New(1, tc.alg)
				hf.AddTarget(sampleHash)
				result = hf.Check([]byte("test"))
				if result == "" {
					t.Errorf("Check() = '', want success")
				}

				// mess up 4 bytes
				sampleHash[1] ^= 0xff
				sampleHash[2] ^= 0xff
				sampleHash[3] ^= 0xff
				hf = New(4, tc.alg)
				hf.AddTarget(sampleHash)
				result = hf.Check([]byte("test"))
				if result == "" {
					t.Errorf("Check() = '', want success")
				}
				hf = New(0, tc.alg)
				hf.AddTarget(sampleHash)
				result = hf.Check([]byte("test"))
				if result != "" {
					t.Errorf("Check() = %q, want ''", result)
				}
			})
		})
	}
}

func BenchmarkHashFinder(b *testing.B) {
	for _, tc := range []struct {
		name string
		algs []func() hash.Hash
	}{
		{
			"SHA1",
			[]func() hash.Hash{sha1.New},
		},
		{
			"SHA256",
			[]func() hash.Hash{sha256.New},
		},
		{
			"MD4",
			[]func() hash.Hash{md4.New},
		},
		{
			"MD5",
			[]func() hash.Hash{md5.New},
		},
		{
			"RipeMD160",
			[]func() hash.Hash{ripemd160.New},
		},
		{
			"5Algs",
			[]func() hash.Hash{
				sha1.New,
				sha256.New,
				md4.New,
				md5.New,
				ripemd160.New,
			},
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.StopTimer()
			// take the 999th hash of a string that seems similar to what we're looking for.
			const text = "Jerry deserves a raise."
			sample := tc.algs[len(tc.algs)-1]()
			sample.Write([]byte(text))
			sampleHash := sample.Sum(nil)
			for i := 0; i < 512; i++ {
				sample.Reset()
				sample.Write(sampleHash)
				sampleHash = sample.Sum(nil)
			}
			// tickle some random bits just to make it interesting
			sampleHash[0] ^= 0x80
			sampleHash[1] ^= 0xff
			sampleHash[2] ^= 0x01
			sampleHash[3] ^= 0x03
			// make a HashFinder with tolerance 4 and set it to look for the sample
			hf := New(4, tc.algs...)
			hf.AddTarget(sampleHash)
			b.StartTimer()

			for i := 0; i < b.N; i++ {
				if result := hf.Check([]byte(text)); result == "" {
					b.Fatalf("Check() = '', want success")
				}
			}
		})
	}
}
