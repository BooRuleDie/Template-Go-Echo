package auth

import (
	constant "go-echo-template/internal/shared/constant"
	"go-echo-template/internal/shared/response"

	"github.com/labstack/echo/v4"
)

const UserContextKey constant.ContextKey = "user"

// GetUserFromContext retrieves the user from the echo context
func GetUserFromContext(c echo.Context) (*User, bool) {
	user, ok := c.Get(string(UserContextKey)).(*User)
	return user, ok
}

// Middleware for authentication
func (s *sessionCookie) CheckAuth(isOptional bool, roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := s.check(c)
			if err != nil {
				if isOptional {
					return next(c)
				}
				return err
			}

			// If specific roles are required, check if the user has one of them
			if len(roles) > 0 {
				hasRole := false
				for _, requiredRole := range roles {
					if user.Role == requiredRole {
						hasRole = true
						break
					}
				}
				if !hasRole {
					return response.ErrSessionUnauthorized
				}
			}

			return next(c)
		}
	}
}
