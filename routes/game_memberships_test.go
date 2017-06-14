package routes

import "fmt"
import "bytes"
import "testing"
import "net/url"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/defs"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testing/utils"
import "github.com/dadleyy/charcoal.api/testing/routing"

func Test_Routes_GameMemberships(t *testing.T) {
	db := testutils.NewDB()
	client, user := models.Client{}, models.User{Email: "game-members@test.com"}

	testutils.CreateClient(&client, "game-memberships-client", false)
	defer db.Unscoped().Delete(&client)

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	game := models.Game{
		Name:    "membership-test-game",
		Status:  "ACTIVE",
		OwnerID: user.ID,
	}
	db.Create(&game)
	defer db.Unscoped().Delete(&game)

	t.Run("update membership", func(update *testing.T) {
		update.Run("update valid status w/ valid membership", func(sub *testing.T) {
			membership := models.GameMembership{UserID: user.ID, GameID: game.ID, Status: "ACTIVE"}
			db.Create(&membership)
			defer db.Unscoped().Delete(&membership)

			json := fmt.Sprintf("{\"status\": \"%s\"}", defs.GameMembershipInactiveStatus)
			reader, params := bytes.NewReader([]byte(json)), testrouting.TestRouteParams{Values: make(url.Values)}
			params.Set("id", fmt.Sprintf("%d", membership.ID))
			ctx := testrouting.NewPatch(&params, reader)
			ctx.Request.User = user
			go func() { <-ctx.Streams["games"] }()
			r := UpdateGameMembership(ctx.Request)

			if r != nil {
				sub.Fatalf("expected nil response but received %v", r)
			}

			var m models.GameMembership
			db.Where("game_id = ? AND user_id = ?", game.ID, user.ID).First(&m)

			if m.Status != defs.GameMembershipInactiveStatus {
				sub.Fatalf("expected %s status but received %v", defs.GameMembershipInactiveStatus, m)
			}
		})

		update.Run("updating another members status", func(other *testing.T) {
			other.Run("update valid status as admin", func(sub *testing.T) {
				player := models.User{Email: "game-members-other@sizethree.cc"}
				db.Create(&player)
				defer db.Unscoped().Delete(&player)

				roleMapping := models.UserRoleMapping{UserID: user.ID, RoleID: 1}
				db.Create(&roleMapping)
				defer db.Unscoped().Delete(&roleMapping)

				membership := models.GameMembership{UserID: player.ID, GameID: game.ID, Status: "ACTIVE"}
				db.Create(&membership)
				defer db.Unscoped().Delete(&membership)

				json := fmt.Sprintf("{\"status\": \"%s\"}", defs.GameMembershipInactiveStatus)
				reader, params := bytes.NewReader([]byte(json)), testrouting.TestRouteParams{Values: make(url.Values)}
				params.Set("id", fmt.Sprintf("%d", membership.ID))

				ctx := testrouting.NewPatch(&params, reader)
				ctx.Request.User = user
				go func() { <-ctx.Streams["games"] }()
				r := UpdateGameMembership(ctx.Request)

				if r != nil {
					sub.Fatalf("expected nil response but received %v", r)
				}

				var m models.GameMembership
				db.Where("game_id = ? AND user_id = ?", game.ID, player.ID).First(&m)

				if m.Status != defs.GameMembershipInactiveStatus {
					sub.Fatalf("expected %s status but received %v", defs.GameMembershipInactiveStatus, m)
				}
			})

			other.Run("update valid status as non-admin", func(sub *testing.T) {
				player := models.User{Email: "game-members-other@sizethree.cc"}
				db.Create(&player)
				defer db.Unscoped().Delete(&player)

				membership := models.GameMembership{UserID: player.ID, GameID: game.ID, Status: "ACTIVE"}
				db.Create(&membership)
				defer db.Unscoped().Delete(&membership)

				json := fmt.Sprintf("{\"status\": \"%s\"}", defs.GameMembershipInactiveStatus)
				reader, params := bytes.NewReader([]byte(json)), testrouting.TestRouteParams{Values: make(url.Values)}
				params.Set("id", fmt.Sprintf("%d", membership.ID))

				ctx := testrouting.NewPatch(&params, reader)
				ctx.Request.User = user
				r := UpdateGameMembership(ctx.Request)

				if r == nil || len(r.Errors) != 1 {
					sub.Fatalf("expected error response but received %v", r)
				}
			})
		})

		update.Run("update invalid status w/ valid membership", func(sub *testing.T) {
			membership := models.GameMembership{UserID: user.ID, GameID: game.ID, Status: "ACTIVE"}
			db.Create(&membership)
			defer db.Unscoped().Delete(&membership)

			json := fmt.Sprintf("{\"status\": \"%s\"}", "garbage")
			reader, params := bytes.NewReader([]byte(json)), testrouting.TestRouteParams{Values: make(url.Values)}
			params.Set("id", fmt.Sprintf("%d", membership.ID))
			ctx := testrouting.NewPatch(&params, reader)
			ctx.Request.User = user
			go func() { <-ctx.Streams["games"] }()
			r := UpdateGameMembership(ctx.Request)

			if r == nil {
				sub.Fatalf("expected error response but received nil")
			}

			if len(r.Errors) != 1 {
				sub.Fatalf("expected error response but received: %v", r)
			}
		})
	})

	t.Run("delete membership", func(del *testing.T) {
		player := models.User{Email: "game-members-delete-test@test.charcoal.sizethree.cc"}
		db.Create(&player)
		defer db.Unscoped().Delete(&player)

		send := func(membershipId uint) *net.ResponseBucket {
			reader, params := bytes.NewReader([]byte("")), testrouting.TestRouteParams{Values: make(url.Values)}
			params.Set("id", fmt.Sprintf("%d", membershipId))

			ctx := testrouting.NewDelete(&params, reader)
			ctx.Request.User = user
			ctx.Request.Client = client
			go func() { <-ctx.Streams["games"] }()
			return DestroyGameMembership(ctx.Request)
		}

		expectSuccess := func(membershipId uint, sub *testing.T) {
			result := send(membershipId)

			if result != nil {
				sub.Fatalf("expected no ResponseBucket but received: %v", result)
				return
			}

			var m models.GameMembership
			e := db.Unscoped().Where("id = ?", membershipId).First(&m).Error

			if e != nil || m.DeletedAt == nil {
				sub.Fatalf("deleted at on membership record was not updated")
			}
		}

		del.Run("delete w/ missing membership id", func(sub *testing.T) {
			reader := bytes.NewReader([]byte(""))
			ctx := testrouting.NewDelete(&testrouting.TestRouteParams{Values: make(url.Values)}, reader)
			ctx.Request.User = user
			ctx.Request.Client = client
			result := DestroyGameMembership(ctx.Request)
			if len(result.Errors) == 0 {
				sub.Fatalf("was expecting error but received: %v", result)
			}
		})

		del.Run("delete w/ invalid membership id", func(sub *testing.T) {
			reader, params := bytes.NewReader([]byte("")), testrouting.TestRouteParams{Values: make(url.Values)}
			params.Set("id", "12313")
			ctx := testrouting.NewDelete(&params, reader)
			ctx.Request.User = user
			ctx.Request.Client = client
			result := DestroyGameMembership(ctx.Request)

			if len(result.Errors) == 0 {
				sub.Fatalf("was expecting error but received: %v", result)
			}
		})

		del.Run("delete by admin", func(sub *testing.T) {
			membership := models.GameMembership{UserID: player.ID, GameID: game.ID, Status: "ACTIVE"}
			db.Create(&membership)
			defer db.Unscoped().Delete(&membership)
			roleMapping := models.UserRoleMapping{UserID: user.ID, RoleID: 1}
			db.Create(&roleMapping)
			defer db.Unscoped().Delete(&roleMapping)
			expectSuccess(membership.ID, sub)
		})

		del.Run("delete by self", func(sub *testing.T) {
			membership := models.GameMembership{UserID: user.ID, GameID: game.ID, Status: "ACTIVE"}
			db.Create(&membership)
			defer db.Unscoped().Delete(&membership)
			expectSuccess(membership.ID, sub)
		})

		del.Run("delete by owner", func(sub *testing.T) {
			membership := models.GameMembership{UserID: player.ID, GameID: game.ID, Status: "ACTIVE"}
			db.Create(&membership)
			defer db.Unscoped().Delete(&membership)
			expectSuccess(membership.ID, sub)
		})

		del.Run("delete by non-self non-owner", func(sub *testing.T) {
			membership := models.GameMembership{UserID: user.ID, GameID: game.ID, Status: "ACTIVE"}
			db.Create(&membership)
			defer db.Unscoped().Delete(&membership)

			reader, params := bytes.NewReader([]byte("")), testrouting.TestRouteParams{Values: make(url.Values)}
			params.Set("id", fmt.Sprintf("%d", membership.ID))

			ctx := testrouting.NewDelete(&params, reader)
			ctx.Request.User = player
			ctx.Request.Client = client
			result := DestroyGameMembership(ctx.Request)

			if len(result.Errors) == 0 {
				sub.Fatalf("was expecting error but received: %v", result)
			}
		})
	})

	t.Run("create membership", func(create *testing.T) {
		player := models.User{Email: "game-members-create-test@test.charcoal.sizethree.cc"}
		db.Create(&player)
		defer db.Unscoped().Delete(&player)

		send := func() *net.ResponseBucket {
			json := fmt.Sprintf("{\"game_id\": %d, \"user_id\": %d}", game.ID, player.ID)
			reader := bytes.NewReader([]byte(json))
			ctx := testrouting.NewPost(&testrouting.TestRouteParams{}, reader)
			ctx.Request.User = user
			ctx.Request.Client = client
			go func() { <-ctx.Streams["games"] }()
			return CreateGameMembership(ctx.Request)
		}

		create.Run("added by owner", func(sub *testing.T) {
			membership := models.GameMembership{UserID: user.ID, GameID: game.ID, Status: "ACTIVE"}
			db.Create(&membership)
			defer db.Unscoped().Delete(&membership)
			result := send()
			defer db.Unscoped().Where("user_id = ? AND game_id = ?", player.ID, game.ID).Delete(models.GameMembership{})
			if len(result.Errors) >= 1 {
				sub.Fatalf("received errors while adding member: %v", result.Errors)
			}
		})

		create.Run("added by admin", func(sub *testing.T) {
			roleMapping := models.UserRoleMapping{UserID: user.ID, RoleID: 1}
			db.Create(&roleMapping)
			defer db.Unscoped().Delete(&roleMapping)
			result := send()
			defer db.Unscoped().Where("user_id = ? AND game_id = ?", player.ID, game.ID).Delete(models.GameMembership{})
			if len(result.Errors) >= 1 {
				sub.Fatalf("received errors while adding member: %v", result.Errors)
			}
		})

		create.Run("added by rando", func(sub *testing.T) {
			result := send()
			if len(result.Errors) == 0 {
				sub.Fatalf("should have received errors")
			}
		})
	})
}
