package eventcount

import (
	"bufio"
	"fmt"
	"io"
)

type Result struct {
	Lines    int64
	NonEmpty int64
	Bytes    int64
}

func Count(reader io.Reader) (Result, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var result Result
	for scanner.Scan() {
		result.Lines++
		result.Bytes += int64(len(scanner.Bytes()))
		if len(scanner.Bytes()) > 0 {
			result.NonEmpty++
		}
	}
	if err := scanner.Err(); err != nil {
		return Result{}, fmt.Errorf("scan events: %w", err)
	}
	return result, nil
}
