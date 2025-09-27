package cache

import (
	"fmt"
	"time"
)

// key prefixes
const (
    userKeyPrefix = "CACHE:USER"
    postKeyPrefix = "CACHE:POST"
)


type cacheKey struct {
	Name string
	TTL  time.Duration
}

// keyToTTL maps cache key prefixes to their TTL (time-to-live) values.
// Keeping all key prefixes centralized in this map helps avoid accidental
// duplication or inconsistent TTL assignments for different cache keys.
var keyToTTL = map[string]time.Duration{
	userKeyPrefix: 24 * time.Hour,
	postKeyPrefix: 24 * time.Hour,
}

// Key builder functions
func GetUserKey(userID int64) *cacheKey {
	ttl, ok := keyToTTL[userKeyPrefix]

	if !ok {
		panic("key must be added to the keyToTTL map")
	}

	return &cacheKey{
		// USER:123
		Name: fmt.Sprintf(userKeyPrefix+":%d", userID),
		TTL:  ttl,
	}
}

func GetPostKey(postID int64) *cacheKey {
	ttl, ok := keyToTTL[postKeyPrefix]

	if !ok {
		panic("key must be added to the keyToTTL map")
	}

	return &cacheKey{
		// POST:123
		Name: fmt.Sprintf(postKeyPrefix+":%d", postID),
		TTL:  ttl,
	}
}
