package net

import "strings"

type Route struct {
	Method string
	Path string
	Handler HandlerFunc
	Middleware []MiddlewareFunc
}

func normalize(path string) []string {
	norm := strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
	return strings.Split(norm, "/")
}

func (route *Route) Reduce() HandlerFunc {
	result := route.Handler

	for _, mw := range route.Middleware {
		result = mw(result)
	}

	return result
}

func (route *Route) Match(method, path string) (UrlParams, bool) {
	params := make(map[string]string)
	result := UrlParams{params}

	if method != route.Method {
		return result, false
	}

	real := normalize(path)
	goal := normalize(route.Path)

	if len(real) != len(goal) {
		return result, false
	}

	for index, part := range goal {
		p := strings.IndexRune(part, ':')

		if p == 0 && len(part) > 1 {
			result.params[part[1:]] = real[index]
			continue
		}

		if part != real[index] {
			return result, false
		}
	}

	return result, true
}
