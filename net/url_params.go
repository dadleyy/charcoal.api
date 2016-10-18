package net

import "strconv"

type UrlParams struct {
	params map[string]string
}

func (p *UrlParams) StringParam(key string) (string, bool) {
	if p.params == nil {
		return "", false
	}

	val, ok := p.params[key]
	return val, ok
}

func (p *UrlParams) IntParam(key string) (int, bool) {
	if p.params == nil {
		return -3, false
	}

	val, ok := p.params[key]

	if ok != true {
		return -2, false
	}

	iv, err := strconv.Atoi(val)

	if err != nil {
		return -1, false
	}

	return iv, true
}
