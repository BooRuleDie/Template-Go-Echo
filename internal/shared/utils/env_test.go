package utils

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// GetStrEnv
func TestGetStrEnv(t *testing.T) {
	t.Run("Existing Environment Variable", func(t *testing.T) {
		// Set up test environment variable
		testValue := "Some Test Value"
		os.Setenv("TEST_STR_ENV", testValue)
		defer os.Unsetenv("TEST_STR_ENV")

		// Test existing environment variable
		result := GetStrEnv("TEST_STR_ENV", "Default")
		require.Equal(t, testValue, result)
	})

	t.Run("Default Value for Non-existent Environment Variable", func(t *testing.T) {
		testDefaultValue := "Some Default Test Value"
		result := GetStrEnv("NON_EXISTENT_ENV", testDefaultValue)
		require.Equal(t, testDefaultValue, result)
	})
}

// MustGetStrEnv
func TestMustGetStrEnv(t *testing.T) {
	t.Run("Existing Environment Variable", func(t *testing.T) {
		testValue := "Some Test Value"
		os.Setenv("TEST_STR_ENV", testValue)
		defer os.Unsetenv("TEST_STR_ENV")

		result := MustGetStrEnv("TEST_STR_ENV")
		require.Equal(t, testValue, result)
	})

	// Test panic on missing environment variable
	t.Run("Panic for Missing Environment Variable", func(t *testing.T) {
		require.Panics(t, func() {
			MustGetStrEnv("NON_EXISTENT_ENV")
		}, "expected panic for missing environment variable")
	})
}

// GetIntEnv
func TestGetIntEnv(t *testing.T) {
	t.Run("Existing Environment Variable", func(t *testing.T) {
		// Set up test environment variable
		testValue := 365
		os.Setenv("TEST_INT_ENV", strconv.Itoa(testValue))
		defer os.Unsetenv("TEST_INT_ENV")

		// Test existing environment variable
		result := GetIntEnv("TEST_INT_ENV", 365)
		require.Equal(t, testValue, result)
	})

	t.Run("Default Value for Non-existent Environment Variable", func(t *testing.T) {
		testDefaultValue := 365
		result := GetIntEnv("NON_EXISTENT_ENV", testDefaultValue)
		require.Equal(t, testDefaultValue, result)
	})

	t.Run("Non-integer Environment Variable", func(t *testing.T) {
		require.Panics(t, func() {
			os.Setenv("TEST_INT_ENV", "Some String Value")
			defer os.Unsetenv("TEST_INT_ENV")

			GetIntEnv("TEST_INT_ENV", 123)
		}, "expected panic for string environment variable")

		require.Panics(t, func() {
			os.Setenv("TEST_INT_ENV", "123.50")
			defer os.Unsetenv("TEST_INT_ENV")

			GetIntEnv("TEST_INT_ENV", 123)
		}, "expected panic for float environment variable")
	})
}

// MustGetIntEnv
func TestMustGetIntEnv(t *testing.T) {
	t.Run("Existing Environment Variable", func(t *testing.T) {
		// Set up test environment variable
		testValue := 365
		os.Setenv("TEST_INT_ENV", strconv.Itoa(testValue))
		defer os.Unsetenv("TEST_INT_ENV")

		// Test existing environment variable
		result := MustGetIntEnv("TEST_INT_ENV")
		require.Equal(t, testValue, result)
	})

	t.Run("Non-integer Environment Variable", func(t *testing.T) {
		require.Panics(t, func() {
			os.Setenv("TEST_INT_ENV", "Some String Value")
			defer os.Unsetenv("TEST_INT_ENV")

			MustGetIntEnv("TEST_INT_ENV")
		}, "expected panic for string environment variable")

		require.Panics(t, func() {
			os.Setenv("TEST_INT_ENV", "123.50")
			defer os.Unsetenv("TEST_INT_ENV")

			MustGetIntEnv("TEST_INT_ENV")
		}, "expected panic for float environment variable")
	})
}

// GetDurationEnv
func TestGetDurationEnv(t *testing.T) {
	t.Run("Existing Environment Variable", func(t *testing.T) {
		testValue := "10s"
		expected := 10 * time.Second
		os.Setenv("TEST_DURATION_ENV", testValue)
		defer os.Unsetenv("TEST_DURATION_ENV")

		result := GetDurationEnv("TEST_DURATION_ENV", 15*time.Second)
		require.Equal(t, expected, result)
	})

	t.Run("Default Value for Non-existent Environment Variable", func(t *testing.T) {
		testDefaultValue := 52 * time.Second
		result := GetDurationEnv("NON_EXISTENT_ENV", testDefaultValue)
		require.Equal(t, testDefaultValue, result)
	})

	t.Run("Invalid Duration Environment Variable", func(t *testing.T) {
		require.Panics(t, func() {
			os.Setenv("TEST_DURATION_ENV", "abc")
			defer os.Unsetenv("TEST_DURATION_ENV")

			GetDurationEnv("TEST_DURATION_ENV", 1*time.Second)
		}, "expected panic for invalid duration value")

		require.Panics(t, func() {
			os.Setenv("TEST_DURATION_ENV", "123") // integer string is not valid duration
			defer os.Unsetenv("TEST_DURATION_ENV")

			GetDurationEnv("TEST_DURATION_ENV", 1*time.Second)
		}, "expected panic for numeric string duration value")
	})
}

// MustGetDurationEnv
func TestMustGetDurationEnv(t *testing.T) {
	t.Run("Existing Environment Variable", func(t *testing.T) {
		testValue := "3m"
		expected := 3 * time.Minute
		os.Setenv("TEST_DURATION_ENV", testValue)
		defer os.Unsetenv("TEST_DURATION_ENV")

		result := MustGetDurationEnv("TEST_DURATION_ENV")
		require.Equal(t, expected, result)
	})

	t.Run("Missing Environment Variable", func(t *testing.T) {
		require.Panics(t, func() {
			MustGetDurationEnv("NON_EXISTENT_ENV")
		}, "expected panic for missing environment variable")
	})

	t.Run("Invalid Duration Environment Variable", func(t *testing.T) {
		require.Panics(t, func() {
			os.Setenv("TEST_DURATION_ENV", "notaduration")
			defer os.Unsetenv("TEST_DURATION_ENV")

			MustGetDurationEnv("TEST_DURATION_ENV")
		}, "expected panic for invalid duration value")

		require.Panics(t, func() {
			os.Setenv("TEST_DURATION_ENV", "42")
			defer os.Unsetenv("TEST_DURATION_ENV")

			MustGetDurationEnv("TEST_DURATION_ENV")
		}, "expected panic for numeric string duration value")
	})
}
