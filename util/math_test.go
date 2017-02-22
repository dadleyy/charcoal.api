package util

import "testing"

func TestMaxInt(t *testing.T) {
	r := MaxInt(1, 100)

	if r != 100 {
		t.Log("expected 100 but received: %d", r)
		t.Fail()
	}
}

func TestMinInt(t *testing.T) {
	r := MinInt(1, 100)

	if r != 1 {
		t.Fail()
	}
}
