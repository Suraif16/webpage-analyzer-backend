package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/suraif16/webpage-analyzer/internal/config"
	"github.com/suraif16/webpage-analyzer/internal/core/domain"
	"go.uber.org/zap"
)

type MockAnalyzer struct {
	mock.Mock
}

func (m *MockAnalyzer) Analyze(ctx context.Context, url string) (*domain.PageAnalysis, error) {
	args := m.Called(ctx, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PageAnalysis), args.Error(1)
}

func TestAnalyzerHandler_Analyze(t *testing.T) {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	config, err := config.LoadConfig()
	if err != nil {
		t.Fatal("Cannot load config:", err)
	}
	gin.SetMode(config.GinMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockAnalyzer)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Successful analysis",
			requestBody: domain.AnalysisRequest{
				URL: "https://example.com",
			},
			setupMock: func(ma *MockAnalyzer) {
				ma.On("Analyze", mock.Anything, "https://example.com").Return(&domain.PageAnalysis{
					HTMLVersion: "HTML5",
					PageTitle:   "Example",
					Headings: domain.HeadingCount{
						H1: 1,
					},
					Links: domain.LinkAnalysis{
						Internal: 2,
						External: 1,
					},
					HasLoginForm: true,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: domain.PageAnalysis{
				HTMLVersion: "HTML5",
				PageTitle:   "Example",
				Headings: domain.HeadingCount{
					H1: 1,
				},
				Links: domain.LinkAnalysis{
					Internal: 2,
					External: 1,
				},
				HasLoginForm: true,
			},
		},
		{
			name: "Invalid URL Format",
			requestBody: map[string]interface{}{
				"url": "invalid-url",
			},
			setupMock: func(ma *MockAnalyzer) {
				// No mock setup needed as error occurs during request binding
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidURL,
		},
		{
			name: "Service Reports Invalid URL",
			requestBody: domain.AnalysisRequest{
				URL: "http://invalid.com",
			},
			setupMock: func(ma *MockAnalyzer) {
				// Mock service-level URL validation failure
				ma.On("Analyze", mock.Anything, "http://invalid.com").
					Return(nil, domain.ErrInvalidURL)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidURL,
		},
		{
			name: "Page not found",
			requestBody: domain.AnalysisRequest{
				URL: "https://example.com/404",
			},
			setupMock: func(ma *MockAnalyzer) {
				// Mock 404 response from service
				ma.On("Analyze", mock.Anything, "https://example.com/404").
					Return(nil, domain.ErrPageNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   domain.ErrPageNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			mockAnalyzer := new(MockAnalyzer)
			tt.setupMock(mockAnalyzer)
			handler := NewAnalyzerHandler(mockAnalyzer, logger)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Analyze(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			expectedJSON, _ := json.Marshal(tt.expectedBody)
			actualJSON, _ := json.Marshal(response)
			assert.JSONEq(t, string(expectedJSON), string(actualJSON))

			mockAnalyzer.AssertExpectations(t)
		})
	}
}
