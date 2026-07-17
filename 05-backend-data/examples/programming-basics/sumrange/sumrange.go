package sumrange

import "fmt"

func SumRange(values []int, start, end int) (int, error) {
	if start < 0 || end > len(values) || start > end {
		return 0, fmt.Errorf("invalid range [%d,%d) for length %d", start, end, len(values))
	}
	total := 0
	for i := start; i < end; i++ {
		total += values[i]
	}
	return total, nil
}
