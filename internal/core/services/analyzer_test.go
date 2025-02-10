package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/suraif16/webpage-analyzer/internal/core/domain"
)

// MockHTTPClient implements the HTTPClient interface for testing
type MockHTTPClient struct {
    mock.Mock
}

func (m *MockHTTPClient) FetchPage(ctx context.Context, url string) (string, error) {
    args := m.Called(ctx, url)
    return args.String(0), args.Error(1)
}

func (m *MockHTTPClient) CheckLink(ctx context.Context, url string) bool {
    args := m.Called(ctx, url)
    return args.Bool(0)
}

// MockHTMLParser implements the HTMLParser interface for testing
type MockHTMLParser struct {
    mock.Mock
}

func (m *MockHTMLParser) GetHTMLVersion(doc string) string {
    args := m.Called(doc)
    return args.String(0)
}

func (m *MockHTMLParser) GetTitle(doc string) string {
    args := m.Called(doc)
    return args.String(0)
}

func (m *MockHTMLParser) CountHeadings(doc string) domain.HeadingCount {
    args := m.Called(doc)
    return args.Get(0).(domain.HeadingCount)
}

func (m *MockHTMLParser) AnalyzeLinks(doc string, baseURL string) domain.LinkAnalysis {
    args := m.Called(doc, baseURL)
    return args.Get(0).(domain.LinkAnalysis)
}

func (m *MockHTMLParser) HasLoginForm(doc string) bool {
    args := m.Called(doc)
    return args.Bool(0)
}

// MockLogger implements the Logger interface for testing
type MockLogger struct {
    mock.Mock
}

func (m *MockLogger) Error(args ...interface{}) {
    m.Called(args...)
}

func (m *MockLogger) Info(args ...interface{}) {
    m.Called(args...)
}

func TestAnalyzerService_Analyze(t *testing.T) {
    // Define test cases
    tests := []struct {
        name          string
        url           string
        setupMocks    func(*MockHTTPClient, *MockHTMLParser, *MockLogger)
        expectedError error
        expectedResult *domain.PageAnalysis
    }{
        {
            name: "Successful analysis",
            url:  "https://example.com",
            setupMocks: func(httpClient *MockHTTPClient, htmlParser *MockHTMLParser, logger *MockLogger) {
                // Setup HTTP client expectations
                httpClient.On("FetchPage", mock.Anything, "https://example.com").
                    Return("<html></html>", nil)

                // Setup HTML parser expectations
                htmlParser.On("GetHTMLVersion", "<html></html>").
                    Return("HTML5")
                htmlParser.On("GetTitle", "<html></html>").
                    Return("Example Title")
                htmlParser.On("CountHeadings", "<html></html>").
                    Return(domain.HeadingCount{H1: 1})
                htmlParser.On("AnalyzeLinks", "<html></html>", "https://example.com").
                    Return(domain.LinkAnalysis{Internal: 1})
                htmlParser.On("HasLoginForm", "<html></html>").
                    Return(false)

                // Setup logger expectations
                logger.On("Info", []interface{}{
                    "page analysis completed",
                    "url",
                    "https://example.com",
                }...).Return()
            },
            expectedError: nil,
            expectedResult: &domain.PageAnalysis{
                HTMLVersion: "HTML5",
                PageTitle:   "Example Title",
                Headings:    domain.HeadingCount{H1: 1},
                Links:       domain.LinkAnalysis{Internal: 1},
                HasLoginForm: false,
            },
        },
        {
            name: "Invalid URL",
            url:  "invalid-url",
            setupMocks: func(httpClient *MockHTTPClient, htmlParser *MockHTMLParser, logger *MockLogger) {
                logger.On("Error", []interface{}{
                    "invalid URL",
                    "url",
                    "invalid-url",
                    "error",
                    mock.AnythingOfType("*url.Error"),
                }...).Return()
            },
            expectedError: domain.ErrInvalidURL,
            expectedResult: nil,
        },
        {
            name: "Page not accessible",
            url:  "https://example.com",
            setupMocks: func(httpClient *MockHTTPClient, htmlParser *MockHTMLParser, logger *MockLogger) {
                httpClient.On("FetchPage", mock.Anything, "https://example.com").
                    Return("", domain.ErrPageNotAccessible)
                logger.On("Error", []interface{}{
                    "failed to fetch page",
                    "url",
                    "https://example.com",
                    "error",
                    domain.ErrPageNotAccessible,
                }...).Return()
            },
            expectedError: domain.ErrPageNotAccessible,
            expectedResult: nil,
        },
    }

    // Execute test cases
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Initialize mocks
            httpClient := new(MockHTTPClient)
            htmlParser := new(MockHTMLParser)
            logger := new(MockLogger)

            // Setup mock expectations
            tt.setupMocks(httpClient, htmlParser, logger)

            // Create service instance
            service := NewAnalyzerService(httpClient, htmlParser, logger)

            // Execute the service method
            result, err := service.Analyze(context.Background(), tt.url)

            // Assert expectations
            if tt.expectedError != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.expectedError, err)
                assert.Nil(t, result)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
                assert.Equal(t, tt.expectedResult, result)
            }

            // Verify all mock expectations were met
            httpClient.AssertExpectations(t)
            htmlParser.AssertExpectations(t)
            logger.AssertExpectations(t)
        })
    }
}