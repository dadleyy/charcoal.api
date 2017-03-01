package net

import "testing"
import "net/url"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/testutils"

func makebp(values url.Values) Blueprint {
	db := testutils.NewDB()
	logger := log.New("miritos")
	return Blueprint{db, logger, values}
}

func Test_Net_Blueprint_LimitUnset(t *testing.T) {
	values := make(url.Values)
	bp := makebp(values)
	l := bp.Limit()

	if l == BlueprintDefaultLimit {
		return
	}

	t.Fatalf("expected default limit but received %d", l)
}

func Test_Net_Blueprint_LimitTooLarge(t *testing.T) {
	values := make(url.Values)
	values.Set("limit", "100000")
	bp := makebp(values)
	l := bp.Limit()

	if l == BlueprintMaxLimit {
		return
	}

	t.Fatalf("expected default limit but received %d", l)
}

func Test_Net_Blueprint_LimitBadParse(t *testing.T) {
	values := make(url.Values)
	values.Set("limit", "abcd")
	bp := makebp(values)
	l := bp.Limit()

	if l == BlueprintDefaultLimit {
		return
	}

	t.Fatalf("expected default limit but received %d", l)
}

func Test_Net_Blueprint_LimitNegative(t *testing.T) {
	values := make(url.Values)
	values.Set("limit", "-100")
	bp := makebp(values)
	l := bp.Limit()

	if l == BlueprintMinLimit {
		return
	}

	t.Fatalf("expected default limit but received %d", l)
}
