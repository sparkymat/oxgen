package api

func BlogsCreate(s internal.Services) echo.HandlerFunc {
  return wrapWithAuth(func(c echo.Context, user dbx.User) error {
    var req Blogs.CreateRequest
    if err := c.Bind(&req); err != nil {
      return c.JSON(http.StatusBadRequest, err)
    }
    res, err := s.Blogs.Create(c.Request().Context(), req)
    if err != nil {
      return c.JSON(http.StatusInternalServerError, err)
    }
    return c.JSON(http.StatusOK, res)
  }
}
