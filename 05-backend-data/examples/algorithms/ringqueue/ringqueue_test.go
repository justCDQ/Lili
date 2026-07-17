package ringqueue

import "testing"

func TestWrapAround(t *testing.T) {
	q, _ := New[string](3)
	for _, value := range []string{"A", "B", "C"} {
		if err := q.Push(value); err != nil {
			t.Fatal(err)
		}
	}
	if got, _ := q.Pop(); got != "A" {
		t.Fatalf("first=%q", got)
	}
	if err := q.Push("D"); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"B", "C", "D"} {
		got, err := q.Pop()
		if err != nil || got != want {
			t.Fatalf("got=%q err=%v want=%q", got, err, want)
		}
	}
}
