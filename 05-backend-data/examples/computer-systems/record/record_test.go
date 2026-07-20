package record

import (
	"encoding/binary"
	"strings"
	"testing"
	"time"
)

func TestRoundTrip(t *testing.T) {
	parsed, err := time.Parse(time.RFC3339, "2026-07-17T10:30:00+08:00")
	if err != nil {
		t.Fatal(err)
	}
	input := Record{Time: parsed, Normalized: true, Text: "狸力"}
	encoded, err := Encode(input)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if !decoded.Time.Equal(input.Time) || decoded.Normalized != input.Normalized || decoded.Text != input.Text {
		t.Fatalf("decoded=%+v input=%+v", decoded, input)
	}
	if got := binary.BigEndian.Uint16(encoded[12:14]); got != 6 {
		t.Fatalf("payload length=%d, want=6", got)
	}
}

func TestDecodeFailures(t *testing.T) {
	valid, err := Encode(Record{Time: time.Unix(0, 1), Text: "ok"})
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		edit func([]byte) []byte
	}{
		{"short header", func([]byte) []byte { return []byte("LI") }},
		{"magic", func(data []byte) []byte { data[0] = 0; return data }},
		{"version", func(data []byte) []byte { data[2] = 2; return data }},
		{"reserved flags", func(data []byte) []byte { data[3] = 0x80; return data }},
		{"length mismatch", func(data []byte) []byte { data[13]++; return data }},
		{"invalid utf8", func(data []byte) []byte { data[14] = 0xff; return data }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := append([]byte(nil), valid...)
			if _, err := Decode(tc.edit(data)); err == nil {
				t.Fatal("Decode() error=nil, want failure")
			}
		})
	}
}

func TestEncodeTooLarge(t *testing.T) {
	if _, err := Encode(Record{Text: strings.Repeat("x", maxPayload+1)}); err == nil {
		t.Fatal("Encode() error=nil, want size failure")
	}
}
