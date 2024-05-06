package route

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Setup(e *echo.Echo, cfg internal.ConfigService, services internal.Services) {
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	app := e.Group("")
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] Got ${method} on ${uri} from ${remote_ip}. Responded with ${status} in ${latency_human}.\n",
	}))
	app.Use(middleware.Recover())

	app.Static("/js", "public/js")
	app.Static("/css", "public/css")
	app.Static("/images", "public/images")
	app.Static("/fonts", "public/fonts")

	// Start of router setup code generated by oxgen. DO NOT EDIT.
	// End of router setup code generated by oxgen. DO NOT EDIT.

	app.Use(session.Middleware(sessions.NewCookieStore([]byte(cfg.SessionSecret()))))

	registerAPIRoutes(app, cfg, services)
}