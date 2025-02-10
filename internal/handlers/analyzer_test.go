package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "context"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/suraif16/webpage-analyzer/internal/core/domain"
)

// MockAnalyzer implements the PageAnalyzer interface for testing
type MockAnalyzer struct {
    mock.Mock
}

// Analyze mocks the analyzer service's Analyze method
func (m *MockAnalyzer) Analyze(ctx context.Context, url string) (*domain.PageAnalysis, error) {
    args := m.Called(ctx, url)
    // If first argument is nil, return nil and the error
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    // Otherwise, return the analysis result and error (which might be nil)
    return args.Get(0).(*domain.PageAnalysis), args.Error(1)
}

func TestAnalyzerHandler_Analyze(t *testing.T) {
    // Set Gin to test mode to avoid debug logging
    gin.SetMode(gin.TestMode)

    // Define test cases using table-driven testing pattern
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
                // Set up expectations for a successful analysis
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

    // Execute each test case
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup test environment
            mockAnalyzer := new(MockAnalyzer)
            tt.setupMock(mockAnalyzer)
            handler := NewAnalyzerHandler(mockAnalyzer)

            // Create test HTTP context
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)

            // Prepare request body
            body, _ := json.Marshal(tt.requestBody)
            c.Request = httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(body))
            c.Request.Header.Set("Content-Type", "application/json")

            // Execute handler
            handler.Analyze(c)

            // Assert HTTP status code
            assert.Equal(t, tt.expectedStatus, w.Code)

            // Parse and verify response body
            var response interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)

            // Compare expected and actual JSON responses
            expectedJSON, _ := json.Marshal(tt.expectedBody)
            actualJSON, _ := json.Marshal(response)
            assert.JSONEq(t, string(expectedJSON), string(actualJSON))

            // Verify all mock expectations were met
            mockAnalyzer.AssertExpectations(t)
        })
    }
}