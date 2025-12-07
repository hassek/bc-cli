package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hassek/bc-cli/config"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Config     *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		BaseURL:    cfg.APIURL,
		HTTPClient: &http.Client{},
		Config:     cfg,
	}
}

func (c *Client) doRequest(method, path string, body any, requireAuth bool) (*http.Response, error) {
	// Check if token needs refresh before making authenticated request
	if requireAuth && c.Config.IsAuthenticated() {
		if c.Config.IsTokenExpired() {
			if c.Config.IsRefreshTokenExpired() {
				return nil, fmt.Errorf("refresh token expired, please login again")
			}

			// Attempt to refresh the token
			if err := c.RefreshToken(); err != nil {
				return nil, fmt.Errorf("failed to refresh token: %w", err)
			}
		}
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.BaseURL + path
	logRequest(method, url, body)

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if requireAuth && c.Config.AccessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Config.AccessToken))
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// If we get 401 Unauthorized and this is an authenticated request, try to refresh and retry once
	if requireAuth && resp.StatusCode == 401 && !c.Config.IsRefreshTokenExpired() {
		_ = resp.Body.Close()

		// Attempt to refresh the token
		if refreshErr := c.RefreshToken(); refreshErr != nil {
			logDebug("Failed to refresh token on 401: %v", refreshErr)
			return resp, nil // Return original 401 response
		}

		// Retry the request with the new token
		if body != nil {
			// Re-marshal the body for the retry
			jsonData, err := json.Marshal(body)
			if err != nil {
				logDebug("Failed to marshal body for retry: %v", err)
				return resp, nil // Return original 401 response
			}
			reqBody = bytes.NewBuffer(jsonData)
		}

		retryReq, err := http.NewRequest(method, url, reqBody)
		if err != nil {
			logDebug("Failed to create retry request: %v", err)
			return resp, nil // Return original 401 response
		}

		retryReq.Header.Set("Content-Type", "application/json")
		retryReq.Header.Set("Accept", "application/json")
		retryReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Config.AccessToken))

		retryResp, err := c.HTTPClient.Do(retryReq)
		if err != nil {
			logDebug("Retry request failed: %v", err)
			return resp, nil // Return original 401 response
		}

		return retryResp, nil
	}

	return resp, nil
}

// logDebug logs debug messages (currently a no-op, but can be enhanced)
func logDebug(format string, args ...any) {
	// In production, this could write to a debug log file
	// For now, we'll keep it silent to avoid cluttering user output
	_ = fmt.Sprintf(format, args...)
}

type APIError struct {
	Data map[string]any `json:"data"`
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Error string `json:"error"`
			Field string `json:"field"`
			Type  string `json:"type"`
		} `json:"errors"`
	} `json:"meta"`
}

func (c *Client) handleResponse(resp *http.Response, result any) error {
	defer func() {
		_ = resp.Body.Close() // Explicitly ignore error in defer
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	logResponse(resp.StatusCode, body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to parse structured API error
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err == nil {
			// Extract field-specific errors
			if len(apiErr.Meta.Errors) > 0 {
				var errorMessages []string
				for _, e := range apiErr.Meta.Errors {
					if e.Field != "" {
						errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", e.Field, e.Error))
					} else {
						errorMessages = append(errorMessages, e.Error)
					}
				}
				if len(errorMessages) > 0 {
					return fmt.Errorf("%s", strings.Join(errorMessages, "\n"))
				}
			}
			// Fallback to meta message
			if apiErr.Meta.Message != "" {
				return fmt.Errorf("%s", apiErr.Meta.Message)
			}
		}

		// Try simple detail message format
		var errResp map[string]any
		if err := json.Unmarshal(body, &errResp); err == nil {
			if msg, ok := errResp["detail"].(string); ok {
				return fmt.Errorf("%s", msg)
			}
		}

		// Fallback to raw response
		return fmt.Errorf("request failed (status %d): %s", resp.StatusCode, string(body))
	}

	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}
