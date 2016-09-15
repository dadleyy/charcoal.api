package main

import "os"
import "fmt"
import "flag"
import "net/http"

import "github.com/labstack/echo"
import "github.com/labstack/echo/engine/standard"

func index(context echo.Context) error {
	return context.String(http.StatusOK, "Hello, World!\n")
}

func main() {
	flag.Parse()
	port := os.Getenv("PORT")

	if len(port) < 1 {
		port = "8080"
	}

	server := echo.New()

	server.Get("/", index)

	server.Run(standard.New(fmt.Sprintf(":%s", port)))
}
