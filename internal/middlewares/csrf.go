package middlewares

import (
	"Go-REST-API-Portfolio/internal/logger"
	"Go-REST-API-Portfolio/internal/utils"
	"crypto/sha256"
	"encoding/base64"
	"github.com/labstack/echo"
	"io"
	"net/http"
)

const (
	// 32 bytes
	csrfSalt = "BbUhoe8qbCC5GEfBa9ovQdzOzXsuVU9I"
)

// Create CSRF token
func MakeToken(sid string, logger logger.Logger) string {
	hash := sha256.New()
	_, err := io.WriteString(hash, csrfSalt+sid)
	if err != nil {
		logger.Errorf("CSRF Token 생성 에러: ", err)
	}
	token := base64.RawStdEncoding.EncodeToString(hash.Sum(nil))
	return token
}

// Validate CSRF token
func ValidateToken(token string, sid string, logger logger.Logger) bool {
	trueToken := MakeToken(sid, logger)
	return token == trueToken
}

func (md *ServerMiddleware) CSRF(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if !md.cfg.Server.CSRF {
			return next(ctx)
		}

		token := ctx.Request().Header.Get("X-CSRF-Token")
		if token == "" {
			md.logger.Errorf("CSRF Token is Empty: token: %s, RequestID: %s", token, utils.GetReqeustID(ctx))
			// echo Context.JSON Function making JSON sends a JSON response with status code
			return ctx.JSON(http.StatusForbidden, utils.NewErrorRes(http.StatusForbidden, "Invaild CSRF Token", "No Token"))
		}

		sid, ok := ctx.Get("sid").(string)
		if !ValidateToken(token, sid, md.logger) || !ok {
			md.logger.Errorf("CSRF Validate Error, token: %s, RequestID: %s", token, utils.GetReqeustID(ctx))
			return ctx.JSON(http.StatusForbidden, utils.NewErrorRes(http.StatusForbidden, "Invaild CSRF Token", "Wrong CSRF Token"))
		}

		return next(ctx)
	}
}