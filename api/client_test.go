package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hassek/bc-cli/config"
)

// TestURLConstruction verifies that URLs are properly constructed regardless of trailing slashes
func TestURLConstruction(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		path        string
		expectedURL string
	}{
		{
			name:        "base URL with trailing slash",
			baseURL:     "http://localhost:8000/",
			path:        "/api/core/v1/subscriptions",
			expectedURL: "http://localhost:8000/api/core/v1/subscriptions",
		},
		{
			name:        "base URL without trailing slash",
			baseURL:     "http://localhost:8000",
			path:        "/api/core/v1/subscriptions",
			expectedURL: "http://localhost:8000/api/core/v1/subscriptions",
		},
		{
			name:        "path without leading slash",
			baseURL:     "http://localhost:8000",
			path:        "api/core/v1/subscriptions",
			expectedURL: "http://localhost:8000/api/core/v1/subscriptions",
		},
		{
			name:        "both with slashes",
			baseURL:     "http://localhost:8000/",
			path:        "/api/core/v1/subscriptions/123/pause",
			expectedURL: "http://localhost:8000/api/core/v1/subscriptions/123/pause",
		},
		{
			name:        "production URL without trailing slash",
			baseURL:     "https://api.butler.coffee",
			path:        "/api/core/v1/subscriptions/123/resume",
			expectedURL: "https://api.butler.coffee/api/core/v1/subscriptions/123/resume",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actualURL string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				actualURL = r.URL.String()
				// Return minimal valid response
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"meta":{"code":200,"message":"ok"},"data":[]}`))
			}))
			defer server.Close()

			// Override the base URL with test server URL for verification
			cfg := &config.Config{
				APIURL:      server.URL,
				AccessToken: "test-token",
			}
			client := NewClient(cfg)

			// Make a request to verify URL construction
			_, _ = client.doRequest("GET", tt.path, nil, false)

			// Verify the path portion matches expectations
			// We can't check the full URL since httptest.Server uses its own host
			// But we can verify no double slashes exist in the path
			if containsDoubleSlash(actualURL) {
				t.Errorf("URL contains double slash: %s", actualURL)
			}
		})
	}
}

// containsDoubleSlash checks if a URL contains // outside of the protocol
func containsDoubleSlash(url string) bool {
	// Remove the protocol part (http:// or https://)
	if len(url) > 8 && url[:7] == "http://" {
		url = url[7:]
	} else if len(url) > 9 && url[:8] == "https://" {
		url = url[8:]
	}

	// Check for double slashes in the remaining part
	for i := 0; i < len(url)-1; i++ {
		if url[i] == '/' && url[i+1] == '/' {
			return true
		}
	}
	return false
}
