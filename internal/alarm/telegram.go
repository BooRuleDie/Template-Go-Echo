package alarm

import (
	"context"
	"fmt"
	"go-echo-template/internal/config"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type telegramAlarmer struct {
	config     config.TelegramConfig
	httpClient *http.Client
}

// Alarm sends a message to the Telegram chat with retry logic.
func (t *telegramAlarmer) Alarm(message string) {
	// Validate input
	if strings.TrimSpace(message) == "" {
		defaultLogger.Printf("ERROR: message cannot be empty")
		return
	}

	apiUrl := "https://api.telegram.org/bot" + t.config.BOT_TOKEN + "/sendMessage"
	values := url.Values{}
	values.Set("chat_id", t.config.CHAT_ID)
	values.Set("text", message)
	body := values.Encode()

	var collectedErrors []error

	for attempt := 1; attempt <= maxRetryCount; attempt++ {
		// fresh nil error for this attempt
		var attemptErr error

		// Create context with timeout for each attempt
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		req, err := http.NewRequestWithContext(ctx, "POST", apiUrl, strings.NewReader(body))
		if err != nil {
			cancel()
			attemptErr = fmt.Errorf("failed to create request: %w", err)
			collectedErrors = append(collectedErrors, attemptErr)
			defaultLogger.Printf("ERROR: %v", attemptErr)
			return
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := t.httpClient.Do(req)
		if err != nil {
			cancel()
			attemptErr = fmt.Errorf("attempt %d failed: %w", attempt, err)
			collectedErrors = append(collectedErrors, attemptErr)
			// Add exponential backoff for retries (except last attempt)
			if attempt < maxRetryCount {
				time.Sleep(time.Duration(attempt) * time.Second)
			}
			continue
		}

		// Process response: we use an anonymous function here to ensure
		// defer statements like resp.Body.Close() and cancel() are always called
		// for proper cleanup, regardless of how the function exits.
		func() {
			defer resp.Body.Close()
			defer cancel()

			if resp.StatusCode != http.StatusOK {
				attemptErr = fmt.Errorf("attempt %d failed: telegram API returned status %s", attempt, resp.Status)
			}
		}()

		if attemptErr != nil {
			collectedErrors = append(collectedErrors, attemptErr)
			// Add backoff for non-successful responses (except last attempt)
			if attempt < maxRetryCount {
				time.Sleep(time.Duration(attempt) * time.Second)
			}
			continue
		} else {
			// success break the loop
			break
		}
	}

	errCount := len(collectedErrors)
	switch errCount {
	case 0:
		// first try success - just return
		return
	case maxRetryCount:
		// all failed - wrap all errors and log
		var errorMessages []string
		for _, err := range collectedErrors {
			errorMessages = append(errorMessages, err.Error())
		}
		wrappedError := fmt.Errorf("all %d attempts failed: %s", maxRetryCount, strings.Join(errorMessages, "; "))
		defaultLogger.Printf("ERROR: %v", wrappedError)
	default:
		// success but errors at the beginning - wrap errors and log
		var errorMessages []string
		for _, err := range collectedErrors {
			errorMessages = append(errorMessages, err.Error())
		}
		wrappedError := fmt.Errorf("succeeded after %d failures: %s", errCount, strings.Join(errorMessages, "; "))
		defaultLogger.Printf("WARNING: %v", wrappedError)
	}
}

// SetGlobalAlarmer initializes and sets the global alarmer with a new default HTTP client.
func SetGlobalAlarmer(cfg *config.TelegramConfig) error {
	if GlobalAlarmer != nil {
		return fmt.Errorf("global alarmer is already set")
	}

	// Validate config
	if strings.TrimSpace(cfg.BOT_TOKEN) == "" {
		return fmt.Errorf("BOT_TOKEN cannot be empty")
	}
	if strings.TrimSpace(cfg.CHAT_ID) == "" {
		return fmt.Errorf("CHAT_ID cannot be empty")
	}

	// Create dedicated HTTP client for this alarmer
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  false,
			MaxIdleConnsPerHost: 2,
		},
	}

	// Create the alarmer
	alarmer := &telegramAlarmer{
		config:     *cfg,
		httpClient: httpClient,
	}

	GlobalAlarmer = alarmer
	defaultLogger.Printf("SUCCESS: Global alarmer initialized successfully")
	return nil
}
