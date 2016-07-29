package users

// import "fmt"
import "github.com/golang/glog"
import "github.com/kataras/iris"
import "github.com/meritoss/meritoss.api/api/models"

func Create(ctx *iris.Context) {
	var user models.User

	if err := ctx.ReadJSON(&user); err != nil {
		glog.Errorf("error reading user: %s\n", err.Error())
		ctx.Panic()
		return
	}

	if len(user.Name) < 2 {
		ctx.SetStatusCode(422)
		ctx.JSON(iris.StatusOK, iris.Map{"error": "bad request"})
		return
	}

	ctx.Write("hi")
}
