package sid

import (
	"testing"
)

func TestNew(t *testing.T) {
	// Generate 10 ids
	ids := make([]ID, 10)
	for i := 0; i < 10; i++ {
		ids[i] = New()
	}

	for i := 1; i < 10; i++ {
		prevID := ids[i-1]
		id := ids[i]

		// Test for uniqueness among all other 9 generated ids
		for j, tid := range ids {
			if j != i {
				if id.Int64() == tid.Int64() {
					t.Errorf("generated ID is not unique (%d/%d)", i, j)
				}
			}
		}

		// Check that timestamp was incremented and is within 30 seconds of the previous one
		secs := id.Time().Sub(prevID.Time()).Seconds()
		if secs < 0 || secs > 30 {
			t.Error("wrong timestamp in generated ID")
		}

		// Check that machine ids are the same
		if id.Node() != prevID.Node() {
			t.Error("machine ID not equal")
		}

		// Test for proper increment
		if got, want := int(id.Step()-prevID.Step()), 1; got != want {
			t.Errorf("wrong increment in generated ID, delta=%v, want %v", got, want)
		}

		t.Logf("id: %d - %d - %d", id.Int64(), id.Node(), id.Step())
	}
}

func BenchmarkNew(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = New()
		}
	})
}

func BenchmarkNewInt64(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = New().Int64()
		}
	})
}
