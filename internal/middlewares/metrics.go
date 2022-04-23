package middlewares

import (
	"Go-REST-API-Portfolio/internal/prometheus"
	"github.com/labstack/echo"
	"time"
)

func (sm *ServerMiddleware) MetricsMiddlewareInit(m prometheus.Metrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()
			err := next(ctx)

			var status int
			if err != nil {
				status = err.(*echo.HTTPError).Code
			} else {
				status = ctx.Response().Status
			}

			/*
				Metrics info injection
			*/
			m.ObeserveResponseTime(status, ctx.Request().Method, ctx.Path(), time.Since(start).Seconds())
			m.IncHits(status, ctx.Request().Method, ctx.Path())
			return err
		}
	}
}