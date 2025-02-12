package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/suraif16/webpage-analyzer/internal/core/domain"
	"go.uber.org/zap"
)

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

func TestAnalyzerService_Analyze(t *testing.T) {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tests := []struct {
		name           string
		url            string
		setupMocks     func(*MockHTTPClient, *MockHTMLParser)
		expectedError  error
		expectedResult *domain.PageAnalysis
	}{
		{
			name: "Successful analysis",
			url:  "https://example.com",
			setupMocks: func(httpClient *MockHTTPClient, htmlParser *MockHTMLParser) {
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
			},
			expectedError: nil,
			expectedResult: &domain.PageAnalysis{
				HTMLVersion:  "HTML5",
				PageTitle:    "Example Title",
				Headings:     domain.HeadingCount{H1: 1},
				Links:        domain.LinkAnalysis{Internal: 1},
				HasLoginForm: false,
			},
		},
		{
			name: "Invalid URL",
			url:  "invalid-url",
			setupMocks: func(httpClient *MockHTTPClient, htmlParser *MockHTMLParser) {
			},
			expectedError:  domain.ErrInvalidURL,
			expectedResult: nil,
		},
		{
			name: "Page not accessible",
			url:  "https://example.com",
			setupMocks: func(httpClient *MockHTTPClient, htmlParser *MockHTMLParser) {
				httpClient.On("FetchPage", mock.Anything, "https://example.com").
					Return("", domain.ErrPageNotAccessible)
			},
			expectedError:  domain.ErrPageNotAccessible,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := new(MockHTTPClient)
			htmlParser := new(MockHTMLParser)

			tt.setupMocks(httpClient, htmlParser)

			service := NewAnalyzerService(httpClient, htmlParser, logger)

			result, err := service.Analyze(context.Background(), tt.url)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult, result)
			}

			httpClient.AssertExpectations(t)
			htmlParser.AssertExpectations(t)
		})
	}
}
