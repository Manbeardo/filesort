package filesort

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"unicode"
	"unicode/utf8"
)

func NoStrAlloc(in io.Reader, out io.Writer) error {
	// Read the input file
	buf, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	lines := [][]byte{}
	i := 0
	for j, b := range buf {
		if b == '\n' {
			lines = append(lines, buf[i:j])
			i = j + 1
		}
	}
	if i < len(buf) {
		lines = append(lines, buf[i:])
	}

	// Trim leading/trailing whitespace
	for i, line := range lines {
		j := 0
		for j < len(line) {
			r, n := utf8.DecodeRune(line[j:])
			if !unicode.IsSpace(r) {
				break
			}
			j += n
		}
		k := len(line)
		for k > j {
			r, n := utf8.DecodeLastRune(line[:k])
			if !unicode.IsSpace(r) {
				break
			}
			k -= n
		}
		lines[i] = line[j:k]
	}

	slices.SortFunc(lines, bytes.Compare)
	lines = slices.CompactFunc(lines, bytes.Equal)

	if len(lines[0]) == 0 {
		lines = lines[1:]
	}

	_, err = out.Write(bytes.Join(lines, []byte("\n")))
	if err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
