package filesort

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"io"
	"runtime"
	"slices"
	"strings"

	"golang.org/x/sync/errgroup"
)

const (
	chunkMaxBytes = L2CacheSize
)

func Streamed(in io.Reader, out io.Writer) error {
	scannedChunkC := make(chan []string)
	processedChunkC := make(chan []string)

	// scanner goroutine
	scannerWG := &errgroup.Group{}
	scannerWG.Go(func() error {
		defer close(scannedChunkC)
		scanner := bufio.NewScanner(in)

		next := []string{}
		nextSize := 0

		for scanner.Scan() {
			lineSize := len(scanner.Bytes())
			if nextSize+lineSize > chunkMaxBytes {
				scannedChunkC <- next
				next = []string{}
				nextSize = 0
			}
			next = append(next, scanner.Text())
			nextSize += len(scanner.Bytes())
		}
		scannedChunkC <- next

		return nil
	})

	// processor goroutines
	processorWG := &errgroup.Group{}
	for range runtime.NumCPU() {
		processorWG.Go(func() error {
			for chunk := range scannedChunkC {
				for i, line := range chunk {
					chunk[i] = strings.TrimSpace(line)
				}
				slices.Sort(chunk)
				chunk = slices.Compact(chunk)
				if chunk[0] == "" {
					chunk = chunk[1:]
				}
				processedChunkC <- chunk
			}
			return nil
		})
	}
	go func() {
		processorWG.Wait()
		close(processedChunkC)
	}()

	// merger/writer goroutine
	writerWG := &errgroup.Group{}
	writerWG.Go(func() error {
		segments := [][]string{{}}
		for chunk := range processedChunkC {
			segments = append(segments, chunk)
			// merge
			for i := len(segments) - 1; i > 0; i-- {
				if len(segments[i]) < len(segments[i-1]) {
					break
				}
				segments[i-1] = mergeSegments(segments[i], segments[i-1])
				segments = segments[:i]
			}
		}
		// final merge
		sorted := []string{}
		for i := len(segments) - 1; i >= 0; i-- {
			sorted = mergeSegments(sorted, segments[i])
		}
		sorted = slices.Compact(sorted)

		// write
		_, err := io.WriteString(out, strings.Join(sorted, "\n"))
		if err != nil {
			return fmt.Errorf("writing output: %w", err)
		}

		return nil
	})

	return errors.Join(
		scannerWG.Wait(),
		processorWG.Wait(),
		writerWG.Wait(),
	)
}

func mergeSegments(a, b []string) []string {
	out := make([]string, len(a)+len(b))
	i, j := 0, 0
	for i+j < len(out) {
		if i >= len(a) {
			out[i+j] = b[j]
			j++
			continue
		}
		if j >= len(b) {
			out[i+j] = a[i]
			i++
			continue
		}
		c := cmp.Compare(a[i], b[j])
		if c == 1 {
			out[i+j] = b[j]
			j++
		} else {
			out[i+j] = a[i]
			i++
		}
	}
	return out
}
