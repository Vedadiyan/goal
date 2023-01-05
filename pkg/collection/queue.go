package collection

import "errors"

type IQueue[T any] interface {
	Enqueue(value T)
	Dequeue() (T, error)
	Peek() (T, error)
	IsEmpty() bool
	Len() int
}

type Queue[T any] struct {
	collection map[int]T
	next       int
	last       int
}

func NewQueue[T any]() Queue[T] {
	queue := Queue[T]{}
	queue.collection = make(map[int]T)
	queue.next = 0
	queue.last = 0
	return queue
}

func (queue *Queue[T]) Enqueue(value T) {
	queue.collection[queue.last] = value
	queue.last++
}

func (queue *Queue[T]) Dequeue() (T, error) {
	if len(queue.collection) == 0 {
		out := new(T)
		return *out, errors.New("queue is empty")
	}
	ref := queue.collection[queue.next]
	delete(queue.collection, queue.next)
	queue.next++
	return ref, nil
}

func (queue *Queue[T]) Peek() (T, error) {
	if len(queue.collection) == 0 {
		out := new(T)
		return *out, errors.New("queue is empty")
	}
	ref := queue.collection[queue.next]
	return ref, nil
}

func (queue *Queue[T]) IsEmpty() bool {
	return len(queue.collection) == 0
}

func (queue *Queue[T]) Len() int {
	return len(queue.collection)
}
