package main

import "fmt"
import "github.com/labstack/echo"
import "github.com/labstack/gommon/log"
import "github.com/labstack/echo/middleware"
import "github.com/labstack/echo/engine/standard"

func Server() *echo.Echo {
	instance := echo.New()

	logger := log.New("miritos")

	logger.SetHeader("[${level}][${short_file}:${line}]")

	instance.SetLogger(logger)

	instance.SetLogLevel(0)

	instance.Use(middleware.Logger())
	instance.Use(middleware.Recover())

	return instance
}

func Run(server *echo.Echo, port string) {
	server.Run(standard.New(fmt.Sprintf(":%s", port)))
}
