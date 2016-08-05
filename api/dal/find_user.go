package dal

import "github.com/golang/glog"
import "github.com/meritoss/meritoss.api/api"
import "github.com/meritoss/meritoss.api/api/middleware"
import "github.com/meritoss/meritoss.api/api/models"

func FindUser(runtime api.Runtime, blueprint middleware.Blueprint) ([]models.User, error) {
	var users []models.User;

	limit, offset := blueprint.Limit, blueprint.Limit * blueprint.Page

	head := runtime.DB.Begin().Limit(limit).Offset(offset)

	for _, filter := range blueprint.Filters {
		glog.Infof("adding filter \"%s\"\n", filter.Reduce())
		head = head.Where(filter.Reduce())
	}

	head.Find(&users).Commit()

	return users, nil
}
