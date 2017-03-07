package net

import "testing"
import "net/url"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/models"
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

func Test_Net_Blueprint_Apply_WithReferencedTable_Matching(t *testing.T) {
	values := url.Values{"filter[game.status]": []string{"eq(ACTIVE)"}}

	bp := makebp(values)

	ownerEmail := "blueprint-reference-test-1@sizethree.cc"
	owner := models.User{Email: ownerEmail}
	game1, game2 := models.Game{Status: "ACTIVE"}, models.Game{Status: "ENDED"}

	bp.Create(&owner)
	defer bp.Unscoped().Delete(&owner)

	game1.OwnerID = owner.ID
	game2.OwnerID = owner.ID

	bp.Create(&game1)
	defer bp.Unscoped().Delete(&game1)

	bp.Create(&game2)
	defer bp.Unscoped().Delete(&game2)

	round := models.GameRound{GameID: game1.ID}
	bp.Create(&round)
	defer bp.Unscoped().Delete(&round)

	var results []models.GameRound

	count, err := bp.Apply(&results)

	if err != nil {
		t.Fatalf("failed: %s", err.Error())
		return
	}

	if count != 1 {
		t.Fatalf("failed: %s", err.Error())
		return
	}
}

func Test_Net_Blueprint_Apply_WithReferencedTable_NoMatch(t *testing.T) {
	values := url.Values{"filter[game.status]": []string{"eq(ENDED)"}}

	bp := makebp(values)

	ownerEmail := "blueprint-reference-test-1@sizethree.cc"
	owner := models.User{Email: ownerEmail}
	game1, game2 := models.Game{Status: "ACTIVE"}, models.Game{Status: "ENDED"}

	bp.Create(&owner)
	defer bp.Unscoped().Delete(&owner)

	game1.OwnerID = owner.ID
	game2.OwnerID = owner.ID

	bp.Create(&game1)
	defer bp.Unscoped().Delete(&game1)

	bp.Create(&game2)
	defer bp.Unscoped().Delete(&game2)

	round := models.GameRound{GameID: game1.ID}
	bp.Create(&round)
	defer bp.Unscoped().Delete(&round)

	var results []models.GameRound

	count, err := bp.Apply(&results)

	if err != nil {
		t.Fatalf("failed: %s", err.Error())
		return
	}

	if count != 0 {
		t.Fatalf("expected 0 results but received %d", count)
		return
	}
}
