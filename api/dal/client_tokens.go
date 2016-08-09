package dal

import "github.com/sizethree/meritoss.api/api/db"
import "github.com/sizethree/meritoss.api/api/models"

func ClientTokensForUser(dbclient *db.Client, user uint) ([]models.ClientToken, int, error) {
	var tokens []models.ClientToken
	var total int

	result := dbclient.Where("user = ?", user).Find(&tokens).Count(&total)

	if e := result.Error; e != nil {
		return tokens, total, e
	}

	return tokens, total, nil
}
