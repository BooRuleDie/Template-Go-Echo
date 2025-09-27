package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"go-echo-template/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

const (
	SessionDefaultExpire = 7 * 24 * time.Hour
	SessionKeyPrefix     = "SESSION:"
	SessionCookieName    = "session"
)

type sessionCookie struct {
	cfg   *config.ServerConfig
	cache *redis.Client
}

// NewSessionCookie creates a new session cookie service
func NewSessionCookie(cfg *config.ServerConfig, cache *redis.Client) AuthService {
	return &sessionCookie{
		cfg:   cfg,
		cache: cache,
	}
}

func (s *sessionCookie) Login(c echo.Context, user *User) error {
	// Generate a secure session ID
	sessionID, err := s.generateSessionID()
	if err != nil {
		return errSessionGenID
	}

	// Create session key with prefix
	sessionKey := SessionKeyPrefix + sessionID

	userJSON, err := json.Marshal(user)
	if err != nil {
		return errSessionSerialize
	}

	// Store session in Redis with expiry
	if err = s.cache.Set(
		c.Request().Context(),
		sessionKey,
		userJSON,
		SessionDefaultExpire,
	).Err(); err != nil {
		return errSessionStore
	}

	// Set session cookie
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

func (s *sessionCookie) Logout(c echo.Context) error {
	// Get session from cookies
	cookie, err := c.Cookie(SessionCookieName)
	if err != nil {
		// Cookie doesn't exist, nothing to logout
		return nil
	}

	sessionID := cookie.Value
	if sessionID == "" {
		return nil
	}

	// Delete session from Redis
	sessionKey := SessionKeyPrefix + sessionID
	if err = s.cache.Del(
		c.Request().Context(),
		sessionKey,
	).Err(); err != nil {
		// errors on logout are not returned
		// nil to .Err() can cause panic, that's why
		// it's wrapped inside that if 
	}

	// Remove session cookie by setting it to expire
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

func (s *sessionCookie) Refresh(c echo.Context, user *User) error {
	// Get session from cookies
	cookie, err := c.Cookie(SessionCookieName)
	if err != nil {
		return errSessionCookieNotFound
	}

	sessionID := cookie.Value
	if sessionID == "" {
		return errEmptySessionID
	}

	// Check if session exists in Redis
	sessionKey := SessionKeyPrefix + sessionID
	exists, err := s.cache.Exists(c.Request().Context(), sessionKey).Result()
	if err != nil {
		return errSessionCheckExist
	}

	if exists == 0 {
		return errSessionNotFound
	}

	if user == nil {
		// Just update the expiry time
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

		// Store updated session in Redis with new expiry
		if err = s.cache.Set(
			c.Request().Context(),
			sessionKey,
			userJSON,
			SessionDefaultExpire,
		).Err(); err != nil {
			return errSessionStore
		}
	}

	// Update cookie expiry
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

func (s *sessionCookie) check(c echo.Context) (*User, error) {
	// Get session from cookies
	cookie, err := c.Cookie(SessionCookieName)
	if err != nil {
		return nil, errSessionCookieNotFound
	}

	sessionID := cookie.Value
	if sessionID == "" {
		return nil, errEmptySessionID
	}

	// Get session data from Redis
	sessionKey := SessionKeyPrefix + sessionID
	userJSON, err := s.cache.Get(c.Request().Context(), sessionKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errSessionNotFound
		}
		return nil, errSessionCheckExist
	}

	// Deserialize user data
	var user User
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return nil, errSessionDeserialize
	}

	// Set user in context for later use
	c.Set(string(UserContextKey), &user)

	return &user, nil
}

// generateSessionID creates a cryptographically secure random session ID
func (s *sessionCookie) generateSessionID() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errSessionGenID
	}
	return hex.EncodeToString(bytes), nil
}
