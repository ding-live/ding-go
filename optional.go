package ding

// Int converts an int to a pointer to an int.
func Int(i int) *int {
	return &i
}

// String converts a string to a pointer to a string.
func String(s string) *string {
	return &s
}
