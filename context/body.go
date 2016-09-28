package context

import "strconv"

type Body map[string]string

func (body *Body) String(key string) (string, bool) {
	result, exists := (*body)[key]
	return result, exists
}

func (body *Body) Int(key string) (int, bool) {
	result, exists := (*body)[key]

	if exists != true {
		return -1, false
	}

	if value, err := strconv.Atoi(result); err == nil {
		return value, true
	}

	return -1, false
}

