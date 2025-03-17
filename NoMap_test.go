package filesort

import (
	"testing"
)

func BenchmarkNoMap(b *testing.B) {
	runBenchmark(b, NoMap)
}

func TestNoMap(t *testing.T) {
	runTest(t, NoMap)
}
