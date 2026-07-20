package counter

import (
	"sync"
	"testing"
)

func TestCounterConcurrent(t *testing.T) {
	counter := New()
	const workers = 8
	const each = 1000
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range each {
				counter.Add("requests")
			}
		}()
	}
	wg.Wait()
	if got, want := counter.Get("requests"), workers*each; got != want {
		t.Fatalf("count=%d, want=%d", got, want)
	}
}

func TestSummary(t *testing.T) {
	counter := New()
	counter.Add("ok")
	if got, want := counter.Summary("ok"), "key=ok count=1"; got != want {
		t.Fatalf("Summary()=%q, want=%q", got, want)
	}
}
