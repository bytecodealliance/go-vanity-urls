package cm

import (
	"unsafe"
)

// List represents a Component Model list<T>.
// The binary representation of list<T> is similar to a Go slice minus the cap field.
type List[T any] struct {
	data *T
	len  uint
}

// NewList returns a List[T] from data and len.
func NewList[T any](data *T, len uint) List[T] {
	return List[T]{
		data: data,
		len:  len,
	}
}

// ToList returns a List[T] equivalent to the Go slice s.
// The underlying slice data is not copied, and the resulting List points at the
// same array storage as the slice.
func ToList[S ~[]T, T any](s S) List[T] {
	return List[T]{
		data: unsafe.SliceData([]T(s)),
		len:  uint(len(s)),
	}
}

// Data returns the data pointer for the list.
func (list List[T]) Data() *T {
	return list.data
}

// Len returns the length of the list.
// TODO: should this return an int instead of a uint?
func (list List[T]) Len() uint {
	return uint(list.len)
}

// Slice returns a Go slice representing the List.
func (list List[T]) Slice() []T {
	return unsafe.Slice(list.data, list.len)
}
