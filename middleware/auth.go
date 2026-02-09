package middleware

import (
	"bambamload/constant"
	"bambamload/enum"
	"bambamload/service/postgresrepository"
	"bambamload/service/redisService"
	"bambamload/utils"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"

	f "github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// GenerateSessionToken generates session token
func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func Authenticate(role enum.SessionOwner, pg postgresrepository.PostgresRepository, redisService redisService.RedisService) f.Handler {
	return func(c *f.Ctx) error {
		authHeader := c.Get(constant.Authorization)
		if authHeader == "" {
			return utils.WriteResponse(c, http.StatusUnauthorized, false, "authorization header is required", nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			return utils.WriteResponse(c, http.StatusUnauthorized, false, "invalid authorization header format", nil)
		}

		token := parts[1]
		redisSession, err := redisService.GetSession(token)
		if err != nil {
			if errors.Is(err, redis.Nil) { //nolint:typecheck
				return utils.WriteResponse(c, http.StatusUnauthorized, false, "invalid credentials. Please login", nil)
			}
			return utils.WriteResponse(c, http.StatusUnauthorized, false, "unable to validate identity, please try again later", nil)
		}

		if redisSession.Owner == enum.SuperAdmin && role == enum.Admin {
			redisSession.Owner = role
		}
		if redisSession.Owner != role {
			return utils.WriteResponse(c, http.StatusUnauthorized, false, "access denied", nil)
		}

		user, err := pg.GetUser(redisSession.ID, constant.ID)
		if err != nil {
			return utils.WriteResponse(c, http.StatusUnauthorized, false, "failed to retrieve user information", nil)
		}

		//if !user.IsLoggedIn {
		//	return utils.WriteResponse(c, http.StatusUnauthorized, false, "unauthorized, please log in", nil)
		//}

		// Store token and user in Fiber context
		c.Locals("token", token)
		c.Locals("user", user)

		return c.Next()
	}
}
