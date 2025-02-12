package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/suraif16/webpage-analyzer/internal/core/domain"
	"github.com/suraif16/webpage-analyzer/internal/core/ports"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// @title Web Page Analyzer API
// @version 1.0
// @description API for analyzing web pages

// @host localhost:8080
// @BasePath /
type AnalyzerHandler struct {
	analyzer ports.PageAnalyzer
	logger   *zap.Logger
}

func NewAnalyzerHandler(analyzer ports.PageAnalyzer, logger *zap.Logger) *AnalyzerHandler {
	return &AnalyzerHandler{
		analyzer: analyzer,
		logger:   logger,
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
	startTime := time.Now()

	var req domain.AnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.ErrInvalidURL)
		return
	}

	h.logger.Info("analyzing url", zap.String("url", req.URL))

	analysis, err := h.analyzer.Analyze(c.Request.Context(), req.URL)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			h.logger.Error("analysis failed",
				zap.String("url", req.URL),
				zap.Error(apiErr))
			c.JSON(apiErr.StatusCode, apiErr)
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrInternalServer)
		return
	}

	h.logger.Info("analysis completed successfully",
		zap.String("url", req.URL),
		zap.String("title", analysis.PageTitle),
		zap.String("html_version", analysis.HTMLVersion),
		zap.Int("headings_h1", analysis.Headings.H1),
		zap.Int("headings_h2", analysis.Headings.H2),
		zap.Int("headings_h3", analysis.Headings.H3),
		zap.Int("headings_h4", analysis.Headings.H4),
		zap.Int("headings_h5", analysis.Headings.H5),
		zap.Int("headings_h6", analysis.Headings.H6),
		zap.Int("internal_links", analysis.Links.Internal),
		zap.Int("external_links", analysis.Links.External),
		zap.Int("inaccessible_links", analysis.Links.Inaccessible),
		zap.Bool("has_login_form", analysis.HasLoginForm),
		zap.Duration("duration", time.Since(startTime)),
	)
	c.JSON(http.StatusOK, analysis)
}
