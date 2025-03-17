package filesort

import (
	"testing"
)

func BenchmarkNoStrAlloc(b *testing.B) {
	runBenchmark(b, NoStrAlloc)
}

func TestNoStrAlloc(t *testing.T) {
	runTest(t, NoStrAlloc)
}
