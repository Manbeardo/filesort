package filesort

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Cache sizes of the R5 4600H
const (
	L1CacheSize = 64 * 1024
	L2CacheSize = 512 * 1024
	L3CacheSize = 8 * 1024 * 1024
)

//go:embed rockyou.txt
var TestFileString string

func createTestFile() (string, error) {
	f, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}
	defer f.Close()
	_, err = f.WriteString(TestFileString)
	if err != nil {
		return "", fmt.Errorf("writing to temp file: %w", err)
	}
	return f.Name(), nil
}

func runBenchmark(b *testing.B, fn func(io.Reader, io.Writer) error) {
	path, err := createTestFile()
	require.NoError(b, err)
	defer os.RemoveAll(path)

	for b.Loop() {
		in, err := os.Open(path)
		require.NoError(b, err)
		out, err := os.CreateTemp(os.TempDir(), "")
		require.NoError(b, err)
		b.ResetTimer()

		err = fn(in, out)

		b.StopTimer()
		in.Close()
		out.Close()
		os.RemoveAll(out.Name())
		if err != nil {
			b.Error(err)
		}
	}
}

var naiveOutputOnce = sync.OnceValue(func() string {
	out := &bytes.Buffer{}
	err := Naive(strings.NewReader(TestFileString), out)
	if err != nil {
		panic(err)
	}
	return out.String()
})

func runTest(t *testing.T, fn func(io.Reader, io.Writer) error) {
	t.Parallel()

	fnOut := &bytes.Buffer{}

	err := fn(strings.NewReader(TestFileString), fnOut)
	require.NoError(t, err)

	assert.Equal(t, naiveOutputOnce(), fnOut.String())
}
