package main

// indexOf returns the index of an element in an array
func indexOf[T comparable](collection []T, e T) int {
	for i, x := range collection {
		if x == e {
			return i
		}
	}
	return -1
}
