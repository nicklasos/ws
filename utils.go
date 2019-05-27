package main

func inArray(value string, array []string) bool {
	for _, val := range array {
		if val == value {
			return true
		}
	}

	return false
}
