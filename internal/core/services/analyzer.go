package services

import (
	"context"
	"net/url"

	"github.com/suraif16/webpage-analyzer/internal/core/domain"
	"github.com/suraif16/webpage-analyzer/internal/core/ports"
	"go.uber.org/zap"
)

type analyzerService struct {
	httpClient ports.HTTPClient
	htmlParser ports.HTMLParser
	logger   *zap.Logger
}


func NewAnalyzerService(httpClient ports.HTTPClient, htmlParser ports.HTMLParser, logger *zap.Logger) ports.PageAnalyzer {
	return &analyzerService{
		httpClient: httpClient,
		htmlParser: htmlParser,
		logger:     logger,
	}
}

func (s *analyzerService) Analyze(ctx context.Context, urlStr string) (*domain.PageAnalysis, error) {
	
	// Validate URL
    if _, err := url.ParseRequestURI(urlStr); err != nil {
        s.logger.Error("invalid URL",
            zap.String("url", urlStr),
            zap.Error(err))
        return nil, domain.ErrInvalidURL
    }

	// Fetch page content
    content, err := s.httpClient.FetchPage(ctx, urlStr)
    if err != nil {
        s.logger.Error("failed to fetch page",
            zap.String("url", urlStr),
            zap.Error(err))
        switch {
        case err == context.DeadlineExceeded:
            return nil, domain.ErrTimeout
        default:
            return nil, domain.ErrPageNotAccessible
        }
    }

	// Analyze page
	s.logger.Info("parsing webpage content")
	analysis := &domain.PageAnalysis{
		HTMLVersion:  s.htmlParser.GetHTMLVersion(content),
		PageTitle:    s.htmlParser.GetTitle(content),
		Headings:     s.htmlParser.CountHeadings(content),
		Links:        s.htmlParser.AnalyzeLinks(content, urlStr),
		HasLoginForm: s.htmlParser.HasLoginForm(content),
	}

    s.logger.Info("page analysis completed",
        zap.String("url", urlStr))
		
	return analysis, nil
}
