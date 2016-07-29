package users

//import "fmt"
import "github.com/golang/glog"
import "github.com/jinzhu/gorm"
import "github.com/kataras/iris"
import "github.com/meritoss/meritoss.api/api"
import "github.com/meritoss/meritoss.api/api/models"


func Index(ctx *iris.Context) {
	runtime, ok := ctx.Get("runtime").(api.Runtime)

	if !ok {
		ctx.Panic()
		return
	}

	var users []models.User;
	result := runtime.DB.Where("name = ?", "jinzhu").Find(&users)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		ctx.JSON(iris.StatusOK, iris.Map{"results": users})
		return
	}

	if e := result.Error; e != nil {
		glog.Errorf("error finding user: %s\n", e.Error())
		ctx.Panic()
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{"results": users})
}
