package topo

import (
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	got, err := Sort([]string{"build", "deploy", "lint", "test"}, []Edge{{"lint", "build"}, {"test", "build"}, {"build", "deploy"}})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"lint", "test", "build", "deploy"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestCycle(t *testing.T) {
	if _, err := Sort([]string{"a", "b"}, []Edge{{"a", "b"}, {"b", "a"}}); err == nil {
		t.Fatal("Sort() error=nil, want cycle")
	}
}
