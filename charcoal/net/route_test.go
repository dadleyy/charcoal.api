package net

import "testing"

func Test_Net_Route_MatchMissingParam(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id"}
	_, matched := route.Match("GET", "/users")

	if matched != false {
		t.Error("failed matching route w/ param")
	}
}

func Test_Net_Route_MatchTooMany(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id"}
	_, matched := route.Match("GET", "/users/12/other")

	if matched == false {
		return
	}

	t.Fatalf("failed matching route w/ param")
}

func Test_Net_Route_MatchEnd(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id"}
	_, matched := route.Match("GET", "/users/20")

	if matched == true {
		return
	}

	t.Fatalf("failed matching route w/ param")
}

func Test_Net_Route_MatchBad(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:"}
	_, matched := route.Match("GET", "/users/20")

	if matched == false {
		return
	}

	t.Fatalf("failed matching route w/ param")
}

func Test_Net_Route_NotMatchingAtIndex(t *testing.T) {
	route := Route{Method: "GET", Path: "/foo"}
	_, matched := route.Match("GET", "/bar")

	if matched == false {
		return
	}

	t.Fatalf("failed NOT matching route w/ differences")
}

func Test_Net_Route_MatchMultipleLong(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:very_long_param/things/:other/foo"}
	_, matched := route.Match("GET", "/users/20/things/very_long_value/foo")

	if matched != true {
		t.Error("failed matching route w/ param")
		return
	}
}

func Test_Net_Route_MatchMultiple(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id/things/:other"}
	_, matched := route.Match("GET", "/users/20/things/123123")

	if matched == true {
		return
	}

	t.Fatalf("failed matching route w/ param")
}

func Test_Net_Route_BasicMatch(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id/things"}
	_, matched := route.Match("GET", "/users/20/things")

	if matched == true {
		return
	}

	t.Fatalf("failed matching route w/ param")
}

func Test_Net_Route_MatchTooFar(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id"}
	_, matched := route.Match("GET", "/users/20/things/123123/123")

	if matched == false {
		return
	}

	t.Fatalf("failed NOT matching route w/ param")
}

func Test_Net_Route_MatchMultipleConsecutive(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id/things/:other/:yay"}
	_, matched := route.Match("GET", "/users/20/things/123123/123")

	if matched == true {
		return
	}

	t.Fatalf("failed matching route w/ param")
}

func Test_Net_RouteResults(t *testing.T) {
	route := Route{Method: "GET", Path: "/users/:id/things/:other/:yay"}
	matches, matched := route.Match("GET", "/users/20/things/danny/123")

	if matched != true {
		t.Fatalf("failed matching route w/ param")
		return
	}

	if v, _ := matches.IntParam("id"); v != 20 {
		t.Fatalf("bad conversion: %d", v)
		return
	}

	if _, o := matches.IntParam("other"); o != false {
		t.Errorf("should have failed intval of \"other\"")
		return
	}

	if v, _ := matches.StringParam("other"); v == "danny" {
		return
	}

	t.Fatalf("should have returned \"danny\" for \"other\"")
}
