package record

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"
)

const headerSize = 14
const maxPayload = 4096

type Record struct {
	Time       time.Time
	Normalized bool
	Text       string
}

func Encode(record Record) ([]byte, error) {
	payload := []byte(record.Text)
	if !utf8.Valid(payload) {
		return nil, errors.New("payload is not valid UTF-8")
	}
	if len(payload) > maxPayload {
		return nil, fmt.Errorf("payload is %d bytes, maximum is %d", len(payload), maxPayload)
	}
	out := make([]byte, headerSize+len(payload))
	copy(out[0:2], []byte{'L', 'I'})
	out[2] = 1
	if record.Normalized {
		out[3] = 1
	}
	binary.BigEndian.PutUint64(out[4:12], uint64(record.Time.UnixNano()))
	binary.BigEndian.PutUint16(out[12:14], uint16(len(payload)))
	copy(out[14:], payload)
	return out, nil
}

func Decode(data []byte) (Record, error) {
	if len(data) < headerSize {
		return Record{}, errors.New("record is shorter than header")
	}
	if data[0] != 'L' || data[1] != 'I' {
		return Record{}, errors.New("invalid magic")
	}
	if data[2] != 1 {
		return Record{}, fmt.Errorf("unsupported version %d", data[2])
	}
	if data[3]&^byte(1) != 0 {
		return Record{}, fmt.Errorf("reserved flag bits set: 0x%02x", data[3])
	}
	size := int(binary.BigEndian.Uint16(data[12:14]))
	if size > maxPayload {
		return Record{}, fmt.Errorf("payload length %d exceeds maximum", size)
	}
	if len(data) != headerSize+size {
		return Record{}, fmt.Errorf("length mismatch: header=%d actual=%d", size, len(data)-headerSize)
	}
	payload := data[14:]
	if !utf8.Valid(payload) {
		return Record{}, errors.New("payload is not valid UTF-8")
	}
	nanos := int64(binary.BigEndian.Uint64(data[4:12]))
	return Record{
		Time:       time.Unix(0, nanos).UTC(),
		Normalized: data[3]&1 != 0,
		Text:       string(payload),
	}, nil
}
