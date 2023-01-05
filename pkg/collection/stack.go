package collection

import (
	"errors"
)

type IStack[T any] interface {
	Push(value T)
	Pop() (T, error)
	Peek() (T, error)
	IsEmpty() bool
	Len() int
}

type Stack[T any] struct {
	collection map[int]T
}

func NewStack[T any]() Stack[T] {
	stack := Stack[T]{}
	stack.collection = make(map[int]T)
	return stack
}

func (stack *Stack[T]) Push(value T) {
	len := len(stack.collection)
	stack.collection[len] = value
}

func (stack *Stack[T]) Pop() (T, error) {
	len := len(stack.collection)
	if len == 0 {
		out := new(T)
		return *out, errors.New("stack is empty")
	}
	ref := stack.collection[len-1]
	delete(stack.collection, len-1)
	return ref, nil
}
func (stack *Stack[T]) Peek() (T, error) {
	len := len(stack.collection)
	if len == 0 {
		out := new(T)
		return *out, errors.New("stack is empty")
	}
	ref := stack.collection[len-1]
	return ref, nil
}

func (stack *Stack[T]) IsEmpty() bool {
	return len(stack.collection) == 0
}

func (stack *Stack[T]) Len() int {
	return len(stack.collection)
}
