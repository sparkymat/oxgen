package route

import "github.com/labstack/echo/v4"

func registerAPIRoutes(app *echo.Group, cfg internal.ConfigService, services internal.Services) {
	apiGroup := app.Group("api")
}
