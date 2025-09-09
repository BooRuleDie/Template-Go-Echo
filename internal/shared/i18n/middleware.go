package i18n

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const localeCookieName = "locale"

// LocaleMiddleware adds locale information to the echo.Context for each request
func LocaleMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(localeCookieName)
		var locale Locale
		if err == nil && cookie != nil {
			// Use existing cookie
			locale = validateLocale(Locale(cookie.Value))
		} else {
			// Set default locale cookie if not existed
			locale = DefaultLocale
			cookie := &http.Cookie{
				Name:     localeCookieName,
				Value:    string(locale),
				Path:     "/",
				HttpOnly: false, // There is no harm JS to access to that cookie
			}
			c.SetCookie(cookie)
		}
		c.Set(localeCookieName, locale)
		return next(c)
	}
}

// Get locale from echo context
func GetLocaleFromContext(c echo.Context) Locale {
	localeStr := c.Get(localeCookieName)
	if localeStr == nil {
		return DefaultLocale
	}

	switch v := localeStr.(type) {
	case Locale:
		return validateLocale(v)
	case string:
		return validateLocale(Locale(v))
	default:
		return DefaultLocale
	}
}

// Make sure only allowed locales are returned
func validateLocale(locale Locale) Locale {
	switch locale {
	case EN_US, TR_TR:
		return locale
	default:
		return DefaultLocale
	}
}
