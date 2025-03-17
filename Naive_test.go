package filesort

import (
	"testing"
)

func BenchmarkNaive(b *testing.B) {
	runBenchmark(b, Naive)
}
