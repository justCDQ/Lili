package workerpool

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Job struct{ ID, Text string }
type Result struct{ ID, Output string }
type indexedJob struct {
	index int
	job   Job
}
type indexedResult struct {
	index  int
	result Result
	err    error
}

func process(ctx context.Context, job Job) (Result, error) {
	select {
	case <-ctx.Done():
		return Result{}, ctx.Err()
	default:
	}
	if job.ID == "" {
		return Result{}, errors.New("job id is empty")
	}
	if job.Text == "FAIL" {
		return Result{}, fmt.Errorf("job %s failed", job.ID)
	}
	return Result{ID: job.ID, Output: strings.ToUpper(job.Text)}, nil
}

func Run(ctx context.Context, jobs []Job, workers int) ([]Result, error) {
	if workers < 1 {
		return nil, errors.New("workers must be positive")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	jobCh := make(chan indexedJob, workers)
	resultCh := make(chan indexedResult, len(jobs))
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobCh {
				result, err := process(ctx, item.job)
				resultCh <- indexedResult{item.index, result, err}
				if err != nil {
					cancel()
				}
			}
		}()
	}
	go func() {
		defer close(jobCh)
		for index, job := range jobs {
			select {
			case jobCh <- indexedJob{index, job}:
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() { wg.Wait(); close(resultCh) }()
	results := make([]Result, len(jobs))
	completed := 0
	var firstErr error
	for item := range resultCh {
		completed++
		if item.err != nil && firstErr == nil {
			firstErr = item.err
		}
		if item.err == nil {
			results[item.index] = item.result
		}
	}
	if firstErr != nil {
		return nil, firstErr
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if completed != len(jobs) {
		return nil, fmt.Errorf("completed %d of %d jobs", completed, len(jobs))
	}
	return results, nil
}
