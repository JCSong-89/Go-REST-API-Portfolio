package middlewares

import (
	"Go-REST-API-Portfolio/config"
	"Go-REST-API-Portfolio/internal/logger"
)

type ServerMiddleware struct {
	//sesstionService
	//authService
	cfg     *config.Config
	origins []string
	logger  logger.Logger
}

func NewMiddleware(
//ss sesstionService,
//as authService,
	c *config.Config,
	origins []string,
	logger logger.Logger) *ServerMiddleware {
	return &ServerMiddleware{
		c, origins, logger,
	}
}