package eventcount

import (
	"bytes"
	"strings"
	"testing"
)

func TestCount(t *testing.T) {
	got, err := Count(strings.NewReader("start\n\nready\n"))
	if err != nil {
		t.Fatal(err)
	}
	want := Result{Lines: 3, NonEmpty: 2, Bytes: 10}
	if got != want {
		t.Fatalf("got=%+v want=%+v", got, want)
	}
}

func TestCountRejectsLongLine(t *testing.T) {
	if _, err := Count(strings.NewReader(strings.Repeat("x", 1024*1024+1))); err == nil {
		t.Fatal("Count() error=nil, want long-line failure")
	}
}

func BenchmarkCount(b *testing.B) {
	data := bytes.Repeat([]byte("event payload\n"), 10000)
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for range b.N {
		got, err := Count(bytes.NewReader(data))
		if err != nil || got.Lines != 10000 {
			b.Fatalf("got=%+v err=%v", got, err)
		}
	}
}
