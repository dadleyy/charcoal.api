package util

import "testing"

func TestMinInt(t *testing.T) {
	r := MaxInt(1, 100)

	if r != 100 {
		t.Fail()
	}
}

func TestMaxInt(t *testing.T) {
	r := MaxInt(1, 100)

	if r != 1 {
		t.Fail()
	}
}
