package sliceutils

import (
	"math"
	"reflect"
	"sort"

	"golang.org/x/exp/constraints"
)

// check if the given slice contains the given elemement.
func Contains[T comparable](s []T, elem T) bool {
	m := make(map[T]bool)
	for _, k := range s {
		m[k] = true
	}

	return m[elem]
}

func Find[T any](s []T, callbackFn func(elem T) bool) (result *T) {
	for _, item := range s {
		if callbackFn(item) {
			return &item
		}
	}

	return nil
}

// returns a new slice containing the given slice's elements, which pass the provided callback function ('callbackFn' returns true).
func Filter[T any](s []T, callbackFn func(elem T) bool) (result []T) {
	for _, item := range s {
		if !callbackFn(item) {
			continue
		}
		result = append(result, item)
	}

	return
}

// returns a new slice containing the given slice's elements after being transformed by the provided callback function.
func Map[A, B any](input []A, callbackFn func(a A) B) (result []B) {
	for _, a := range input {
		b := callbackFn(a)
		result = append(result, b)
	}

	return
}

// returns true if at lease one element of the input array passes the implementation of provided callback function
// otherwise return false.
func Some[T any](s []T, callbackFn func(elem T) bool) bool {
	for _, item := range s {
		if callbackFn(item) {
			return true
		}
	}

	return false
}

func IndexOf[T comparable](s []T, elem T) int {
	for i, item := range s {
		if item == elem {
			return i
		}
	}

	return -1
}

// Check if 2 slices have the same elements without considering their order
// Supported types: int, float, string.
func Equal[T constraints.Ordered](a []T, b []T) bool {
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i] < b[j]
	})
	return reflect.DeepEqual(a, b)
}

// returns the first element of the given slice, which pass the provided callback function ('callbackFn' returns true).

func IsEmpty[T any](s []T) bool {
	return s == nil || len(s) == 0
}

func Reduce[A, B any](input []A, callbackFn func(accVal B, currentVal A) B, initVal B) B {
	accumulation := initVal
	for _, a := range input {
		current := a
		accumulation = callbackFn(accumulation, current)
	}

	return accumulation
}

// return a new slice of unique elements based on the given slice.
func Uniquify[T comparable](s []T) (result []T) {
	m := make(map[T]reflect.Value)
	for _, k := range s {
		v := reflect.ValueOf(k)
		m[k] = v
	}

	for k, v := range m {
		if v.String() == "" {
			continue
		}
		result = append(result, k)
	}

	return
}

// divide the given slice into equal parts whose max length is the given 'maxSize'.
func DivideIntoChunks[T any](s []T, maxSize int64) (result [][]T) {
	parts := int64(math.Ceil(float64(len(s)) / float64(maxSize)))

	for i := int64(0); i < parts; i++ {
		start := i * maxSize
		end := math.Min(float64((i+1)*maxSize), float64(len(s)))

		result = append(result, s[start:int(end)])
	}

	return
}

func FlatMap[A, B any](input []A, f func(A) []B) []B {
	var result []B
	for _, v := range input {
		result = append(result, f(v)...)
	}
	return result
}

func UniqueByField[T any](items []T, field string) []T {
	uniqueValues := make(map[interface{}]struct{})
	var result []T

	for _, item := range items {
		v := reflect.ValueOf(item)
		f := reflect.Indirect(v).FieldByName(field)
		if f.IsValid() {
			if _, exists := uniqueValues[f.Interface()]; !exists {
				uniqueValues[f.Interface()] = struct{}{}
				result = append(result, item)
			}
		}
	}

	return result
}
