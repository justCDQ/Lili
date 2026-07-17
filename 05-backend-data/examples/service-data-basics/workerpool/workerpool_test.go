package workerpool

import (
	"context"
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	jobs := []Job{{"1", "go"}, {"2", "http"}, {"3", "sql"}}
	got, err := Run(context.Background(), jobs, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := []Result{{"1", "GO"}, {"2", "HTTP"}, {"3", "SQL"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestRunFailure(t *testing.T) {
	if _, err := Run(context.Background(), []Job{{"1", "FAIL"}}, 1); err == nil {
		t.Fatal("Run() error=nil, want failure")
	}
}
