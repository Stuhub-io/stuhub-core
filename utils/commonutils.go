package commonutils

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

func GetSlugResolution(existingSlugs []string, slug string) string {
	originalSlug := slug
	counter := 1

	// Create a map for faster lookup
	existingSlugMap := make(map[string]bool)
	for _, s := range existingSlugs {
		existingSlugMap[s] = true
	}

	// Increment the slug until a unique one is found
	for {
		if _, exists := existingSlugMap[slug]; !exists {
			break
		}

		slug = originalSlug + "-" + strconv.Itoa(counter)
		counter++
	}

	return slug
}

func RemoveDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := make([]T, 0, len(sliceList))
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func ToLowerCaseStringList(list []string) []string {
	for i, v := range list {
		list[i] = strings.ToLower(v)
	}
	return list
}

func Filter[T any](slice []T, test func(T) bool) []T {
	result := []T{}
	for _, v := range slice {
		if test(v) {
			result = append(result, v)
		}
	}
	return result
}

func IsIntegerType(field reflect.Value) bool {
	// Check for signed integers
	if field.Kind() == reflect.Int || field.Kind() == reflect.Int8 || field.Kind() == reflect.Int16 ||
		field.Kind() == reflect.Int32 || field.Kind() == reflect.Int64 {
		return true
	}
	// Check for unsigned integers
	if field.Kind() == reflect.Uint || field.Kind() == reflect.Uint8 || field.Kind() == reflect.Uint16 ||
		field.Kind() == reflect.Uint32 || field.Kind() == reflect.Uint64 {
		return true
	}
	return false
}

func CurTimestampAsFloat64() float64 {
	return float64(time.Now().UnixNano()) / 1000000000
}

func NillableField[T any](value T) *T {
	return &value
}
