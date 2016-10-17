package net


type Multiplexer struct {
	routes []Route
	middleware []MiddlewareFunc
}

// wrap
//
// Given a handler function, wrap will apply all of the global middleware onto it
// and return the wrapped handler function.
func (mux *Multiplexer) wrap(handler HandlerFunc) HandlerFunc {
	for _, mw := range mux.middleware {
		handler = mw(handler)
	}

	return handler
}

// Find
//
// Given a method and a path, this function will iterate over all known routes,
// checking to see if it is able to match the current route. If so, it will then
// apply all of the global middleware onto the reduced route.
func (mux *Multiplexer) Find(method, path string) (HandlerFunc, UrlParams, bool) {
	var noop HandlerFunc

	for _, route := range mux.routes {
		// check to see if this one matches
		if params, match := route.Match(method, path); match == true {
			return mux.wrap(route.Reduce()), params, match
		}
	}

	return noop, UrlParams{}, false
}

// routing functions

func (mux *Multiplexer) add(method string, path string, handler HandlerFunc, middleware []MiddlewareFunc) {
	mux.routes = append(mux.routes, Route{method, path, handler, middleware})
}

func (mux *Multiplexer) Use(mw MiddlewareFunc) {
	mux.middleware = append(mux.middleware, mw)
}

func (mux *Multiplexer) GET(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	mux.add("GET", path, handler, middleware)
}

func (mux *Multiplexer) POST(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	mux.add("POST", path, handler, middleware)
}
