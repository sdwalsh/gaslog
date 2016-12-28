package main

import (
	"html/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:_token",
	}))

}
