package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	// Enable logs
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", Index)
	e.GET("/initialize", Initialize)
	e.GET("/list", List)
	e.PUT("/insert", Insert)
	e.DELETE("/id/:deleteid", Deletebyid)
	// This one is used by Puppet fact
	e.GET("/role/fulltext/:hostname", Getrolebyhostname)

	// Start server
	e.Logger.Fatal(e.Start(":14002"))
}
