package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/suraif16/webpage-analyzer/internal/core/domain"
	"github.com/suraif16/webpage-analyzer/internal/core/services"
	"github.com/suraif16/webpage-analyzer/internal/handlers"
	httpclient "github.com/suraif16/webpage-analyzer/internal/infrastructure/http/client" // Updated import alias
	"github.com/suraif16/webpage-analyzer/internal/infrastructure/parser"
	"go.uber.org/zap"
)

// setupRouter creates a test router with all necessary dependencies
func setupRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()

    // Initialize logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    // Initialize dependencies with proper error handling
    httpClient := httpclient.NewHTTPClient(10 * time.Second)
    htmlParser := parser.NewHTMLParser()
    analyzerService := services.NewAnalyzerService(httpClient, htmlParser, logger.Sugar())
    handler := handlers.NewAnalyzerHandler(analyzerService)

    r.POST("/analyze", handler.Analyze)
    return r
}

func TestIntegrationAnalyze(t *testing.T) {
    router := setupRouter()

    tests := []struct {
        name           string
        requestBody    interface{}
        expectedStatus int
        expectedError  *domain.APIError
    }{
        {
            name: "Valid URL",
            requestBody: domain.AnalysisRequest{
                URL: "https://www.example.com",
            },
            expectedStatus: nethttp.StatusOK,
            expectedError:  nil,
        },
        {
            name: "Invalid URL Format",
            requestBody: domain.AnalysisRequest{
                URL: "not-a-url",
            },
            expectedStatus: nethttp.StatusBadRequest,
            expectedError:  domain.ErrInvalidURL,
        },
        {
            name: "Malformed Request",
            requestBody: map[string]interface{}{
                "invalid_field": "value",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  domain.ErrInvalidURL,
        },
        {
            name: "Empty URL",
            requestBody: domain.AnalysisRequest{
                URL: "",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  domain.ErrInvalidURL,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Prepare request
            body, err := json.Marshal(tt.requestBody)
            assert.NoError(t, err, "Failed to marshal request body")

            req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedStatus, w.Code)

            // Parse and verify response
            if tt.expectedError != nil {
                var response domain.APIError
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.NoError(t, err, "Failed to unmarshal error response")
                assert.Equal(t, tt.expectedError.StatusCode, response.StatusCode)
                assert.Equal(t, tt.expectedError.Message, response.Message)
            } else if w.Code == nethttp.StatusOK {
                var response domain.PageAnalysis
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.NoError(t, err, "Failed to unmarshal success response")
                // Verify basic structure of response
                assert.NotEmpty(t, response.HTMLVersion)
                assert.NotEmpty(t, response.PageTitle)
            }
        })
    }
}