package server

import "fmt"
import "github.com/labstack/echo"
import "github.com/labstack/gommon/log"
import "github.com/labstack/echo/middleware"
import "github.com/labstack/echo/engine/standard"

const logfmt = "[${level} ${prefix} ${short_file}:${line}]"

type App struct {
	*echo.Echo
}

func NewApp() *App {
	instance := &App{echo.New()}

	logger := log.New("miritos")
	logger.SetHeader(logfmt)
	instance.SetLogger(logger)
	instance.SetLogLevel(0)

	instance.Use(middleware.Recover())

	return instance
}

func RunApp(app *App, port string) {
	app.Run(standard.New(fmt.Sprintf(":%s", port)))
}
