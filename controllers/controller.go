package controllers

import "github.com/labstack/echo/v4"

func HandleError(c echo.Context, err error, code int) error {
	return c.JSON(code, map[string]interface{}{
		"error": err.Error(),
	})
}

func HandleSuccess(c echo.Context, data interface{}, code int) error {
	return c.JSON(code, map[string]interface{}{
		"payload": data,
	})
}
