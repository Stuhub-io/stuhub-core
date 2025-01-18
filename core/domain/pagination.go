package domain

import (
	"fmt"
	"reflect"
)

const (
	SmallPageSize  = 10
	MediumPageSize = 20
	LargePageSize  = 50
	SuperLargeSize = 100
	GiantPageSize  = 200
)

type Pagination struct {
	Size int64 `json:"size"`
	Page int64 `json:"page"`
}

type OffsetBasedPagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type Cursor interface{ comparable }

type CursorPagination[K Cursor] struct {
	Cursor     K   `json:"cursor,omitempty"`
	NextCursor K   `json:"next_cursor"`
	Limit      int `json:"limit,omitempty"`
}

func CalculateNextCursor[T any, K Cursor](limit int, items []T, fieldName string) *K {
	if len(items) == 0 || len(items) < limit {
		return nil
	}

	val := reflect.ValueOf(items)
	if val.Kind() != reflect.Slice {
		fmt.Println("Expected a slice, but got:", val.Kind())
		return nil
	}

	lastItem := items[len(items)-1]
	lastItemVal := reflect.ValueOf(lastItem)
	if lastItemVal.Kind() == reflect.Struct {
		fieldValue := lastItemVal.FieldByName(fieldName)
		if fieldValue.IsValid() {
			fmt.Println(">>>> ", fieldValue.Interface())
			if fieldValue.Kind() == reflect.Ptr {
				return fieldValue.Interface().(*K)
			} else {
				valCopy := fieldValue.Interface()
				result := valCopy.(K)
				return &result
			}
		}
	}

	return nil
}
