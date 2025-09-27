package auth

import (
	"time"

	"github.com/labstack/echo/v4"
)

type AuthService interface {
	Login(c echo.Context, user *User) error
	Logout(c echo.Context) error
	Refresh(c echo.Context, user *User) error
	check(c echo.Context) (*User, error)

	// middleware
	CheckAuth(isOptional bool, roles ...string) echo.MiddlewareFunc
}

type User struct {
	ID        int64
	Name      string
	Email     string
	Phone     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
