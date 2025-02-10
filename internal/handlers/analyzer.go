package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/suraif16/webpage-analyzer/internal/core/domain"
    "github.com/suraif16/webpage-analyzer/internal/core/ports"
)

type AnalyzerHandler struct {
    analyzer ports.PageAnalyzer
}

func NewAnalyzerHandler(analyzer ports.PageAnalyzer) *AnalyzerHandler {
    return &AnalyzerHandler{
        analyzer: analyzer,
    }
}

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