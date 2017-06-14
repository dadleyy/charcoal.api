package routes

import "fmt"
import "bytes"
import "testing"
import "github.com/dadleyy/charcoal.api/net"
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

	t.Run("delete membership", func(create *testing.T) {
	})

	t.Run("create membership", func(create *testing.T) {
		player := models.User{Email: "game-members-2@test.com"}
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

		t.Run("added by owner", func(sub *testing.T) {
			membership := models.GameMembership{UserID: user.ID, GameID: game.ID, Status: "ACTIVE"}
			db.Create(&membership)
			defer db.Unscoped().Delete(&membership)
			result := send()
			defer db.Unscoped().Where("user_id = ? AND game_id = ?", player.ID, game.ID).Delete(models.GameMembership{})
			if len(result.Errors) >= 1 {
				sub.Fatalf("received errors while adding member: %v", result.Errors)
			}
		})

		t.Run("added by admin", func(sub *testing.T) {
			roleMapping := models.UserRoleMapping{UserID: user.ID, RoleID: 1}
			db.Create(&roleMapping)
			defer db.Unscoped().Delete(&roleMapping)
			result := send()
			defer db.Unscoped().Where("user_id = ? AND game_id = ?", player.ID, game.ID).Delete(models.GameMembership{})
			if len(result.Errors) >= 1 {
				sub.Fatalf("received errors while adding member: %v", result.Errors)
			}
		})

		t.Run("added by rando", func(sub *testing.T) {
			result := send()
			if len(result.Errors) == 0 {
				sub.Fatalf("should have received errors")
			}
		})
	})
}
