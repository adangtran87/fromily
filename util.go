package main

func InStringArray(s string, array []string) bool {
	if array == nil {
		return false
	}
	for _, array_str := range array {
		if s == array_str {
			return true
		}
	}
	return false
}
