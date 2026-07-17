package topk

import (
	"reflect"
	"testing"
)

func TestSelect(t *testing.T) {
	input := []Item{{"a", 90}, {"b", 70}, {"c", 90}, {"d", 80}, {"e", 95}}
	got, err := Select(input, 3)
	if err != nil {
		t.Fatal(err)
	}
	want := []Item{{"e", 95}, {"a", 90}, {"c", 90}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}
