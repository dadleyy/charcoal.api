package dal

import "github.com/golang/glog"
import "github.com/meritoss/meritoss.api/api"
import "github.com/meritoss/meritoss.api/api/middleware"
import "github.com/meritoss/meritoss.api/api/models"

func FindUser(runtime api.Runtime, blueprint middleware.Blueprint) ([]models.User, int, error) {
	var users []models.User
	var total int

	limit, offset := blueprint.Limit, blueprint.Limit * blueprint.Page

	head := runtime.DB.Begin().Limit(limit).Offset(offset)

	for _, filter := range blueprint.Filters {
		glog.Infof("adding filter \"%s\"\n", filter.Reduce())
		head = head.Where(filter.Reduce())
	}

	result := head.Find(&users).Count(&total).Commit()

	if result.Error != nil {
		glog.Errorf("error in FindUser: %s\n", result.Error.Error())
		return users, -1, result.Error
	}

	return users, total, nil
}
