package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/suraif16/webpage-analyzer/internal/core/domain"
    "github.com/suraif16/webpage-analyzer/internal/core/ports"
)

// @title Web Page Analyzer API
// @version 1.0
// @description API for analyzing web pages

// @host localhost:8080
// @BasePath /
type AnalyzerHandler struct {
    analyzer ports.PageAnalyzer
}

func NewAnalyzerHandler(analyzer ports.PageAnalyzer) *AnalyzerHandler {
    return &AnalyzerHandler{
        analyzer: analyzer,
    }
}

// Analyze godoc
// @Summary Analyze a webpage
// @Description Analyzes a webpage for HTML version, headings, links, and login form
// @Tags analyzer
// @Accept json
// @Produce json
// @Param request body domain.AnalysisRequest true "URL to analyze"
// @Success 200 {object} domain.PageAnalysis
// @Failure 400 {object} domain.APIError
// @Failure 404 {object} domain.APIError
// @Failure 500 {object} domain.APIError
// @Router /analyze [post]
func (h *AnalyzerHandler) Analyze(c *gin.Context) {
    var req domain.AnalysisRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, domain.ErrInvalidURL)
        return
    }

    analysis, err := h.analyzer.Analyze(c.Request.Context(), req.URL)
    if err != nil {
        if apiErr, ok := err.(*domain.APIError); ok {
            c.JSON(apiErr.StatusCode, apiErr)
            return
        }
        c.JSON(http.StatusInternalServerError, domain.ErrInternalServer)
        return
    }

    c.JSON(http.StatusOK, analysis)
}