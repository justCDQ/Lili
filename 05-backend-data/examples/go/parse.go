package examplesgo

import (
	"fmt"
	"strconv"
	"strings"
)

func ParsePositiveInt(input string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, fmt.Errorf("parse positive integer %q: %w", input, err)
	}
	if value < 1 {
		return 0, fmt.Errorf("parse positive integer %q: value must be positive", input)
	}
	return value, nil
}
