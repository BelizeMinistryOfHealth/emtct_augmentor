package strs

// Index returns the index position of a string in an array of strings.
// Returns -1 if the string does not exist in the array.
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Include indicates if a string is present in an array.
func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}
