package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/suraif16/webpage-analyzer/internal/core/domain"
	"go.uber.org/zap"
)

func TestHTTPClient_FetchPage(t *testing.T) {

	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedError  error
	}{
		{
			name: "Successful request",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("<html><body>Hello</body></html>"))
			},
			expectedError: nil,
		},
		{
			name: "404 response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectedError: domain.ErrPageNotFound,
		},
		{
			name: "500 response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedError: &domain.APIError{StatusCode: 500},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			client := NewHTTPClient(5 * time.Second, logger)
			content, err := client.FetchPage(context.Background(), server.URL)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if apiErr, ok := err.(*domain.APIError); ok {
					assert.Equal(t, tt.expectedError.(*domain.APIError).StatusCode, apiErr.StatusCode)
				} else {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Contains(t, content, "Hello")
			}
		})
	}
}

func TestHTTPClient_CheckLink(t *testing.T) {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expected       bool
	}{
		{
			name: "Valid link",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expected: true,
		},
		{
			name: "Invalid link",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			client := NewHTTPClient(5 * time.Second, logger)
			result := client.CheckLink(context.Background(), server.URL)
			assert.Equal(t, tt.expected, result)
		})
	}
}
