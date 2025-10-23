package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"go-echo-template/internal/config"
	"go-echo-template/internal/shared"
	"go-echo-template/internal/shared/log"
	"go-echo-template/internal/shared/utils"
	"go-echo-template/internal/storage"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

const (
	SessionDefaultExpire = 7 * 24 * time.Hour
	SessionKeyPrefix     = "SESSION:"
	SessionCookieName    = "session"

	UserContextKey shared.ContextKey = "user"
)

type AuthService interface {
	// Generic cookie-based session management
	Login(c echo.Context, user *User) error
	Logout(c echo.Context) error
	Refresh(c echo.Context, user *User) error
	Check(c echo.Context) (*User, error)

	// Middleware for general session enforcement
	CheckAuth(isOptional bool, roles ...string) echo.MiddlewareFunc

	// Auth API methods (handler specific)
	apiLogin(c echo.Context, req *LoginRequest) error
	apiRefresh(c echo.Context) error
}

// Session user data
type User struct {
	ID        int64
	Name      string
	Email     string
	Phone     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type service struct {
	cfg     *config.ServerConfig
	cache   *redis.Client
	logger  log.CustomLogger
	storage *storage.Storage
}

func NewSessionCookieService(
	cfg *config.ServerConfig,
	logger log.CustomLogger,
	cache *redis.Client,
	storage *storage.Storage,
) AuthService {
	return &service{logger: logger, storage: storage, cache: cache, cfg: cfg}
}

// --- GENERIC SESSION METHODS ---

// Login sets the session in cookie and redis
func (s *service) Login(c echo.Context, user *User) error {
	// Generate a secure session ID
	sessionID, err := s.generateSessionID()
	if err != nil {
		return errSessionGenID
	}

	sessionKey := SessionKeyPrefix + sessionID

	userJSON, err := json.Marshal(user)
	if err != nil {
		return errSessionSerialize
	}

	if err := s.cache.Set(
		c.Request().Context(),
		sessionKey,
		userJSON,
		SessionDefaultExpire,
	).Err(); err != nil {
		return errSessionStore
	}

	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(SessionDefaultExpire.Seconds()),
		HttpOnly: true,
		Secure:   s.cfg.IsProduction(),
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookie)
	return nil
}

func (s *service) Logout(c echo.Context) error {
	cookie, err := c.Cookie(SessionCookieName)
	if err != nil {
		// If the session cookie does not exist, nothing to do, just return
		return nil
	}

	sessionID := cookie.Value
	if sessionID == "" {
		// If the cookie value is empty, nothing to do, just return
		return nil
	}

	sessionKey := SessionKeyPrefix + sessionID
	if err := s.cache.Del(
		c.Request().Context(),
		sessionKey,
	).Err(); err != nil {
		// just log error if needed
	}

	expiredCookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   s.cfg.IsProduction(),
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(expiredCookie)
	return nil
}

func (s *service) Refresh(c echo.Context, user *User) error {
	cookie, err := c.Cookie(SessionCookieName)
	if err != nil {
		return errSessionCookieNotFound
	}
	sessionID := cookie.Value
	if sessionID == "" {
		return errEmptySessionID
	}

	sessionKey := SessionKeyPrefix + sessionID
	exists, err := s.cache.Exists(c.Request().Context(), sessionKey).Result()
	if err != nil {
		return errSessionCheckExist
	}
	if exists == 0 {
		return errSessionNotFound
	}

	if user == nil {
		if err := s.cache.Expire(
			c.Request().Context(),
			sessionKey,
			SessionDefaultExpire,
		).Err(); err != nil {
			return errSessionStore
		}
	} else {
		userJSON, err := json.Marshal(user)
		if err != nil {
			return errSessionSerialize
		}

		if err := s.cache.Set(
			c.Request().Context(),
			sessionKey,
			userJSON,
			SessionDefaultExpire,
		).Err(); err != nil {
			return errSessionStore
		}
	}

	refreshedCookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(SessionDefaultExpire.Seconds()),
		HttpOnly: true,
		Secure:   s.cfg.IsProduction(),
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(refreshedCookie)
	return nil
}

func (s *service) Check(c echo.Context) (*User, error) {
	cookie, err := c.Cookie(SessionCookieName)
	if err != nil {
		return nil, errSessionCookieNotFound
	}
	sessionID := cookie.Value
	if sessionID == "" {
		return nil, errEmptySessionID
	}

	sessionKey := SessionKeyPrefix + sessionID
	userJSON, err := s.cache.Get(c.Request().Context(), sessionKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errSessionNotFound
		}
		return nil, errSessionCheckExist
	}

	var user User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, errSessionDeserialize
	}
	c.Set(string(UserContextKey), &user)
	return &user, nil
}

// GetUserFromContext retrieves the user from the echo context
func GetUserFromContext(c echo.Context) (*User, bool) {
	user, ok := c.Get(string(UserContextKey)).(*User)
	return user, ok
}

// CheckAuth is authentication middleware, with optional and role-based support
func (s *service) CheckAuth(isOptional bool, roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := s.Check(c)
			if err != nil {
				if isOptional {
					return next(c)
				}
				return shared.ErrSessionUnauthorized
			}

			if len(roles) > 0 {
				hasRole := false
				for _, requiredRole := range roles {
					if user.Role == requiredRole {
						hasRole = true
						break
					}
				}
				if !hasRole {
					return shared.ErrSessionUnauthorized
				}
			}

			return next(c)
		}
	}
}

// --- AUTH API METHODS ---

func (s *service) apiLogin(c echo.Context, req *LoginRequest) error {
	// Get user by email
	userRow, err := s.storage.Auth.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return shared.ErrSessionUnauthorized
	}

	// Check password using bcrypt helper
	if !utils.CheckPasswordHash(req.Password, userRow.Password) {
		return shared.ErrSessionUnauthorized
	}

	user := &User{
		ID:        userRow.ID,
		Name:      userRow.Name,
		Email:     userRow.Email,
		Phone:     userRow.Phone.String,
		Role:      userRow.Role,
		CreatedAt: userRow.CreatedAt,
		UpdatedAt: userRow.UpdatedAt,
	}

	return s.Login(c, user)
}

// APIRefresh: handler-specific auth method for /refresh
func (s *service) apiRefresh(c echo.Context) error {
	// Optionally get associated user (could get from session or DB as needed, here nil to just refresh expiry)
	return s.Refresh(c, nil)
}

// APILogout: handler-specific auth method for /logout
func (s *service) APILogout(c echo.Context) error {
	return s.Logout(c)
}

// generateSessionID creates a cryptographically secure random session ID
func (s *service) generateSessionID() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errSessionGenID
	}
	return hex.EncodeToString(bytes), nil
}
