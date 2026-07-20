package ringqueue

import "errors"

var ErrFull = errors.New("queue is full")
var ErrEmpty = errors.New("queue is empty")

type Queue[T any] struct {
	items []T
	head  int
	size  int
}

func New[T any](capacity int) (*Queue[T], error) {
	if capacity < 1 {
		return nil, errors.New("capacity must be positive")
	}
	return &Queue[T]{items: make([]T, capacity)}, nil
}
func (q *Queue[T]) Len() int { return q.size }
func (q *Queue[T]) Push(value T) error {
	if q.size == len(q.items) {
		return ErrFull
	}
	tail := (q.head + q.size) % len(q.items)
	q.items[tail] = value
	q.size++
	return nil
}
func (q *Queue[T]) Pop() (T, error) {
	if q.size == 0 {
		var zero T
		return zero, ErrEmpty
	}
	value := q.items[q.head]
	var zero T
	q.items[q.head] = zero
	q.head = (q.head + 1) % len(q.items)
	q.size--
	return value, nil
}
