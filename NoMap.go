package filesort

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

func NoMap(in io.Reader, out io.Writer) error {
	// Read the input file
	f, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	// Process the file
	lines := strings.Split(string(f), "\n")

	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	slices.Sort(lines)
	lines = slices.Compact(lines)

	if lines[0] == "" {
		lines = lines[1:]
	}

	_, err = io.WriteString(out, strings.Join(lines, "\n"))
	if err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
