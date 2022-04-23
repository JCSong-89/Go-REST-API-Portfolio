package middlewares

import (
	"Go-REST-API-Portfolio/internal/utils"
	"github.com/labstack/echo"
	"time"
)

func (sm *ServerMiddleware) ReqLogMiddlewareInit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		/*
			Request return latency
		*/
		s := time.Since(start).String()
		requestID := utils.GetReqeustID(ctx)

		sm.logger.Infof("RequestID: %s, Method: %s, URI: %s, Status: %v, Size: %v, Time: %s", requestID, req.Method, req.URL, status, size, s)

		return err
	}
}