package commonutils

import "strconv"

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
