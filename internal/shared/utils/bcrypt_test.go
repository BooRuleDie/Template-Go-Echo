package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	t.Run("Hash and Check Success", func(t *testing.T) {
		password := "mySuperSecret123!"
		hash, err := HashPassword(password)
		require.NoError(t, err, "Hashing should not return error")
		require.NotEmpty(t, hash, "Hash should not be empty")
		require.True(t, CheckPasswordHash(password, hash), "Password should match hash")
	})

	t.Run("Hash and Check Failure (wrong password)", func(t *testing.T) {
		password := "correctPassword"
		wrongPassword := "wrongPassword"
		hash, err := HashPassword(password)
		require.NoError(t, err, "Hashing should not return error")
		require.False(t, CheckPasswordHash(wrongPassword, hash), "Wrong password should not match hash")
	})

	t.Run("CheckPasswordHash with Invalid Hash", func(t *testing.T) {
		invalidHash := "$2a$10$invalidnotavalidhashvalueeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		result := CheckPasswordHash("irrelevant", invalidHash)
		require.False(t, result, "Should return false for invalid hash format")
	})

	t.Run("HashPassword Different Hashes for Same Password", func(t *testing.T) {
		password := "repeatablePassword"
		hash1, err1 := HashPassword(password)
		hash2, err2 := HashPassword(password)
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NotEqual(t, hash1, hash2, "Hashes should be different due to random salt")
		// But both hashes should validate
		require.True(t, CheckPasswordHash(password, hash1))
		require.True(t, CheckPasswordHash(password, hash2))
	})
}
