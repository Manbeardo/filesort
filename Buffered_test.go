package filesort

import (
	"testing"
)

func BenchmarkBuffered(b *testing.B) {
	runBenchmark(b, Buffered)
}

func TestBuffered(t *testing.T) {
	runTest(t, Buffered)
}
