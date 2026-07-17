package examplesgo

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
)

func TestValidateOrder(t *testing.T) {
	err := ValidateOrder("", 0)
	if !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("errors.Is(%v, ErrInvalidOrder) = false", err)
	}
	fieldErr, ok := errors.AsType[*FieldError](err)
	if !ok || fieldErr.Field != "id" {
		t.Fatalf("first FieldError = %#v, %v", fieldErr, ok)
	}
}

func TestSquarePipeline(t *testing.T) {
	got, err := SquarePipeline(t.Context(), []int{2, 3, 4})
	if err != nil {
		t.Fatal(err)
	}
	if want := []int{4, 9, 16}; !slices.Equal(got, want) {
		t.Fatalf("SquarePipeline() = %v, want %v", got, want)
	}
}

func TestSquarePipelineCanceled(t *testing.T) {
	ctx, cancel := context.WithCancelCause(t.Context())
	want := errors.New("caller stopped")
	cancel(want)
	_, err := SquarePipeline(ctx, []int{2, 3, 4})
	if !errors.Is(err, want) {
		t.Fatalf("SquarePipeline() error = %v, want cause %v", err, want)
	}
}

func TestLedgerConcurrentTransfersPreserveTotal(t *testing.T) {
	ledger := NewLedger(map[string]int64{"a": 1000, "b": 1000})
	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			_ = ledger.Transfer("a", "b", 1)
			_ = ledger.Transfer("b", "a", 1)
		})
	}
	wg.Wait()
	if got := ledger.Total(); got != 2000 {
		t.Fatalf("Total() = %d, want 2000", got)
	}
}

func TestRegistry(t *testing.T) {
	var registry Registry
	registry.Publish([]string{"/health", "/orders"})
	for range 100 {
		registry.RecordRequest()
	}
	if got := registry.RequestCount(); got != 100 {
		t.Fatalf("RequestCount() = %d, want 100", got)
	}
	if got := registry.Routes(); !slices.Equal(got, []string{"/health", "/orders"}) {
		t.Fatalf("Routes() = %v", got)
	}
}

func TestRunPool(t *testing.T) {
	jobs := []Job{{ID: 1, Value: 2}, {ID: 2, Value: 3}, {ID: 3, Value: 4}}
	results, err := RunPool(t.Context(), 2, jobs)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != len(jobs) {
		t.Fatalf("len(results) = %d, want %d", len(results), len(jobs))
	}
}

func TestRunPoolCanceled(t *testing.T) {
	ctx, cancel := context.WithCancelCause(t.Context())
	want := errors.New("deployment stopped")
	cancel(want)
	_, err := RunPool(ctx, 2, []Job{{ID: 1, Value: 2}})
	if !errors.Is(err, want) {
		t.Fatalf("RunPool() error = %v, want cause %v", err, want)
	}
}

func TestParsePositiveInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "integer", input: "42", want: 42},
		{name: "trim spaces", input: " 7 ", want: 7},
		{name: "zero", input: "0", wantErr: true},
		{name: "syntax", input: "seven", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParsePositiveInt(test.input)
			if (err != nil) != test.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, test.wantErr)
			}
			if got != test.want {
				t.Fatalf("value = %d, want %d", got, test.want)
			}
		})
	}
}

func BenchmarkParsePositiveInt(b *testing.B) {
	for b.Loop() {
		_, _ = ParsePositiveInt("42")
	}
}
