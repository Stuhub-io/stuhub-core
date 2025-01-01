package commonutils

import (
	"strconv"
	"strings"
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
