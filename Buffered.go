package filesort

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"
)

func Buffered(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	uniqueMap := make(map[string]bool, len(lines))
	var trimmed string
	for _, line := range lines {
		if trimmed = strings.TrimSpace(line); trimmed != "" {
			uniqueMap[trimmed] = true
		}
	}

	// Convert map keys to slice
	ss := make([]string, len(uniqueMap))
	i := 0
	for key := range uniqueMap {
		ss[i] = key
		i++
	}

	slices.Sort(ss)

	_, err := io.WriteString(out, strings.Join(ss, "\n"))
	if err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
