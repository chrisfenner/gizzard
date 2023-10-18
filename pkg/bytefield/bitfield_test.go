package bytefield

import "testing"

func TestBytefield(t *testing.T) {
	var bf Bytefield
	for i := 0; i < 256; i++ {
		if count := bf.CountOf(byte(i)); count != 0 {
			t.Errorf("CountOf(%d) = %d, want 0", i, count)
		}
	}
	bf.Increment(17)
	for i := 0; i < 256; i++ {
		wantCount := 0
		if i == 17 {
			wantCount = 1
		}
		if count := bf.CountOf(byte(i)); count != wantCount {
			t.Errorf("CountOf(%d) = %d, want %d", i, count, wantCount)
		}
	}
	bf.Increment(200)
	bf.Increment(200)
	bf.Increment(200)
	for i := 0; i < 256; i++ {
		wantCount := 0
		if i == 17 {
			wantCount = 1
		} else if i == 200 {
			wantCount = 3
		}
		if count := bf.CountOf(byte(i)); count != wantCount {
			t.Errorf("CountOf(%d) = %d, want %d", i, count, wantCount)
		}
	}
}