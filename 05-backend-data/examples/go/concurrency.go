package examplesgo

import (
	"context"
	"errors"
	"sync"
)

func SquarePipeline(ctx context.Context, input []int) ([]int, error) {
	numbers := make(chan int)
	squares := make(chan int)

	go func() {
		defer close(numbers)
		for _, n := range input {
			select {
			case numbers <- n:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		defer close(squares)
		for n := range numbers {
			select {
			case squares <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()

	var result []int
	for {
		select {
		case value, ok := <-squares:
			if !ok {
				if err := context.Cause(ctx); err != nil {
					return nil, err
				}
				return result, nil
			}
			result = append(result, value)
		case <-ctx.Done():
			return nil, context.Cause(ctx)
		}
	}
}

type Ledger struct {
	mu       sync.Mutex
	balances map[string]int64
}

func NewLedger(initial map[string]int64) *Ledger {
	balances := make(map[string]int64, len(initial))
	for account, amount := range initial {
		balances[account] = amount
	}
	return &Ledger{balances: balances}
}

func (l *Ledger) Transfer(from, to string, amount int64) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if l.balances[from] < amount {
		return errors.New("insufficient balance")
	}
	l.balances[from] -= amount
	l.balances[to] += amount
	return nil
}

func (l *Ledger) Total() int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	var total int64
	for _, amount := range l.balances {
		total += amount
	}
	return total
}

type Job struct {
	ID    int
	Value int
}

type Result struct {
	JobID int
	Value int
}

func RunPool(ctx context.Context, workers int, jobs []Job) ([]Result, error) {
	if workers < 1 {
		return nil, errors.New("workers must be positive")
	}
	jobCh := make(chan Job)
	resultCh := make(chan Result)
	var wg sync.WaitGroup

	for range workers {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobCh:
					if !ok {
						return
					}
					select {
					case resultCh <- Result{JobID: job.ID, Value: job.Value * job.Value}:
					case <-ctx.Done():
						return
					}
				}
			}
		})
	}

	go func() {
		defer close(jobCh)
		for _, job := range jobs {
			select {
			case jobCh <- job:
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	results := make([]Result, 0, len(jobs))
	for result := range resultCh {
		results = append(results, result)
	}
	if err := context.Cause(ctx); err != nil {
		return nil, err
	}
	return results, nil
}
