package net

import "testing"

func TestRouteMatchMissingParam(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id"}
	_, matched := route.Match("GET", "/users")

	if matched != false {
		t.Error("failed matching route w/ param")
	}

	t.Log("successfully didn't match with missing param")
}

func TestRouteMatchTooMany(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id"}
	_, matched := route.Match("GET", "/users/12/other")

	if matched != false {
		t.Error("failed matching route w/ param")
	}

	t.Log("successfully didn't match with missing param")
}

func TestRouteMatchEnd(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id"}
	_, matched := route.Match("GET", "/users/20")

	if matched != true {
		t.Error("failed matching route w/ param")
	}

	t.Log("successfully matched w/ param at end")
}

func TestRouteMatchBad(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:"}
	_, matched := route.Match("GET", "/users/20")

	if matched != false {
		t.Error("failed matching route w/ param")
		return
	}

	t.Log("successfully didn't match with strange ending")
}

func TestNotMatchingAtIndex(t *testing.T) {
	route := Route{Method: "GET", Path: "/foo"}
	_, matched := route.Match("GET", "/bar")

	if matched != false {
		t.Error("failed NOT matching route w/ differences")
		return
	}

	t.Log("successfully matched w/ multiple params")
}

func TestRouteMatchMultipleLong(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:very_long_param/things/:other/foo"}
	_, matched := route.Match("GET", "/users/20/things/very_long_value/foo")

	if matched != true {
		t.Error("failed matching route w/ param")
		return
	}

	t.Log("successfully matched w/ multiple params")
}

func TestRouteMatchMultiple(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id/things/:other"}
	_, matched := route.Match("GET", "/users/20/things/123123")

	if matched != true {
		t.Error("failed matching route w/ param")
	}

	t.Log("successfully matched w/ multiple params")
}

func TestRouteMatch(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id/things"}
	_, matched := route.Match("GET", "/users/20/things")

	if matched != true {
		t.Error("failed matching route w/ param")
	}

	t.Log("successfully matched w/ param in middle")
}

func TestRouteMatchTooFar(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id"}
	_, matched := route.Match("GET", "/users/20/things/123123/123")

	if matched != false {
		t.Error("failed NOT matching route w/ param")
	}

	t.Log("successfully matched w/ multiple consecutive params")
}

func TestRouteMatchMultipleConsecutive(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id/things/:other/:yay"}
	_, matched := route.Match("GET", "/users/20/things/123123/123")

	if matched != true {
		t.Error("failed matching route w/ param")
	}

	t.Log("successfully matched w/ multiple consecutive params")
}

func TestRouteResults(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id/things/:other/:yay"}
	matches, matched := route.Match("GET", "/users/20/things/danny/123")

	if matched != true {
		t.Error("failed matching route w/ param")
	}

	if v, _ := matches.IntParam("id"); v != 20 {
		t.Errorf("bad conversion: %d", v)
		return
	}

	if _, o := matches.IntParam("other"); o != false {
		t.Errorf("should have failed intval of \"other\"")
		return
	}

	if v, _ := matches.StringParam("other"); v != "danny" {
		t.Errorf("should have returned \"danny\" for \"other\"")
		return
	}

	t.Log("successfully parsed url params")
}
