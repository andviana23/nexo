// Package asaas provides a client for the Asaas payment gateway API.
// Documentation: https://docs.asaas.com/reference
package asaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Config holds the Asaas API configuration
type Config struct {
	APIKey      string
	BaseURL     string // https://api.asaas.com/v3 (production) or https://sandbox.asaas.com/api/v3
	Timeout     time.Duration
	MaxRetries  int
	Environment string // "sandbox" or "production"
}

// Client is the Asaas API client
type Client struct {
	httpClient *http.Client
	config     Config
	logger     *zap.Logger
}

// NewClient creates a new Asaas API client
func NewClient(cfg Config, logger *zap.Logger) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}
	if cfg.BaseURL == "" {
		if cfg.Environment == "production" {
			cfg.BaseURL = "https://api.asaas.com/v3"
		} else {
			cfg.BaseURL = "https://sandbox.asaas.com/api/v3"
		}
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		config: cfg,
		logger: logger,
	}
}

// doRequest executes an HTTP request with retry logic
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	url := c.config.BaseURL + path

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt*attempt) * 100 * time.Millisecond
			time.Sleep(backoff)
			c.logger.Debug("retrying request",
				zap.String("url", url),
				zap.Int("attempt", attempt),
			)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("access_token", c.config.APIKey)
		req.Header.Set("User-Agent", "NEXO-Backend/1.0")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			continue
		}

		// Log response
		c.logger.Debug("asaas api response",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", resp.StatusCode),
		)

		// Check for errors
		if resp.StatusCode >= 400 {
			var apiErr ErrorResponse
			if err := json.Unmarshal(respBody, &apiErr); err == nil && len(apiErr.Errors) > 0 {
				return nil, &APIError{
					StatusCode: resp.StatusCode,
					Errors:     apiErr.Errors,
				}
			}
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Errors: []Error{{
					Code:        "UNKNOWN_ERROR",
					Description: string(respBody),
				}},
			}
		}

		return respBody, nil
	}

	return nil, lastErr
}

// ============================================================================
// ERROR TYPES
// ============================================================================

// Error represents a single Asaas API error
type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// ErrorResponse represents the Asaas API error response format
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// APIError represents an error returned by the Asaas API
type APIError struct {
	StatusCode int
	Errors     []Error
}

func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("asaas api error (status %d): %s - %s",
			e.StatusCode, e.Errors[0].Code, e.Errors[0].Description)
	}
	return fmt.Sprintf("asaas api error (status %d)", e.StatusCode)
}

// IsNotFound returns true if the error is a 404 Not Found
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == 404
}

// IsConflict returns true if the error is a 409 Conflict
func (e *APIError) IsConflict() bool {
	return e.StatusCode == 409
}

// IsRateLimited returns true if the error is a 429 Too Many Requests
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == 429
}
