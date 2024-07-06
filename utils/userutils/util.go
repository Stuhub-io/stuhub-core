package userutils

func GetUserFullName(f string, l string) string {
	if f == "" {
		return l
	}

	if l == "" {
		return f
	}

	return f + " " + l
}
