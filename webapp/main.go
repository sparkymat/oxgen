package main

//go:generate go run github.com/valyala/quicktemplate/qtc -dir=internal/view

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sparkymat/oxgen/webapp/internal"
	"github.com/sparkymat/oxgen/webapp/internal/config"
	"github.com/sparkymat/oxgen/webapp/internal/database"
	"github.com/sparkymat/oxgen/webapp/internal/dbx"
	"github.com/sparkymat/oxgen/webapp/internal/route"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	dbDriver, err := database.New(cfg.DatabaseURL())
	if err != nil {
		log.Error(err)
		panic(err)
	}

	if err = dbDriver.AutoMigrate(); err != nil {
		log.Error(err)
		panic(err)
	}

	// Initialize web server
	db := dbx.New(dbDriver.DB())

	services := internal.Services{}

	// Start of main code generated by oxgen. DO NOT EDIT.
	// End of main code generated by oxgen. DO NOT EDIT.

	e := echo.New()
	route.Setup(e, cfg, services)

	e.Logger.Panic(e.Start(":8080"))
}
