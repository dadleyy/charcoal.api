package users

import "github.com/kataras/iris"

func Create(ctx *iris.Context) {
	ctx.Write("Hi %s", "iris")
}
