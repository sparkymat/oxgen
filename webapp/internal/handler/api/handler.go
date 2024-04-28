package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/docker/distribution/registry/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sparkymat/oxgen/webapp/internal/dbx"
)

func renderError(c echo.Context, statusCode int, message string, err error) error {
	if err != nil {
		log.Errorf("err: %v", err)
	}

	//nolint:wrapcheck
	return c.JSON(statusCode, map[string]string{
		"error": message,
	})
}

type (
	authenticatedHandlerFunc       func(c echo.Context, user dbx.User) error
	authenticatedMemberHandlerFunc func(c echo.Context, user dbx.User, id uuid.UUID) error
)

func wrapWithAuth(handlerFunc authenticatedHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, isUser := c.Get(auth.UserKey).(dbx.User)
		if !isUser {
			return renderError(c, http.StatusInternalServerError, "failed to load user", nil)
		}

		return handlerFunc(c, user)
	}
}

func wrapWithAuthForMember(handlerFunc authenticatedMemberHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, isUser := c.Get(auth.UserKey).(dbx.User)
		if !isUser {
			return renderError(c, http.StatusInternalServerError, "failed to load user", nil)
		}

		idString := c.Param("id")

		id, err := uuid.Parse(idString)
		if err != nil {
			return renderError(c, http.StatusNotFound, "not found", err)
		}

		return handlerFunc(c, user, id)
	}
}

//nolint:revive
func parsePaginationParams(c echo.Context) (int32, int32, error) {
	pageSizeString := c.QueryParam("pageSize")

	pageSize, err := strconv.ParseInt(pageSizeString, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("pageSize was invalid. err: %w", err)
	}

	pageNumberString := c.QueryParam("pageNumber")

	pageNumber, err := strconv.ParseInt(pageNumberString, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("pageNumber was invalid. err: %w", err)
	}

	return int32(pageSize), int32(pageNumber), nil
}
