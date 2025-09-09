package i18n

import (
	"fmt"
)

// Strongly typed Locale
type Locale string

const (
	EN_US Locale = "en-US"
	TR_TR Locale = "tr-TR"
)

// Fallback locale
const DefaultLocale = TR_TR

// Translate returns a localized message based on code, locale, and args
func Translate(code string, locale Locale, args ...any) string {
	// Check if code exists
	if locMap, ok := translations[code]; ok {
		// Check if locale exists, else fallback
		if msg, ok := locMap[locale]; ok {
			if len(args) > 0 {
				return fmt.Sprintf(msg, args...)
			}
			return msg
		}
		// fallback locale
		if msg, ok := locMap[DefaultLocale]; ok {
			if len(args) > 0 {
				return fmt.Sprintf(msg, args...)
			}
			return msg
		}
	}
	// ultimate fallback: return code
	return code
}
