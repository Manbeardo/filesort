package filesort

import (
	"testing"
)

func BenchmarkStreamed(b *testing.B) {
	runBenchmark(b, Streamed)
}

func TestStreamed(t *testing.T) {
	runTest(t, Streamed)
}
