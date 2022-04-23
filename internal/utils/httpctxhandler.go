package utils

import "github.com/labstack/echo"

/*
echo Response is wraps http.ResponeWirter
*/
func GetReqeustID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}