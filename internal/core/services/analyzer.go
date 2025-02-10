package services

import (
    "context"
    "net/url"
    "github.com/suraif16/webpage-analyzer/internal/core/domain"
	"github.com/suraif16/webpage-analyzer/internal/core/ports"
)

type analyzerService struct {
    httpClient  ports.HTTPClient
    htmlParser  ports.HTMLParser
    logger      Logger
}

type Logger interface {
    Error(args ...interface{})
    Info(args ...interface{})
}

func NewAnalyzerService(httpClient ports.HTTPClient, htmlParser ports.HTMLParser, logger Logger) ports.PageAnalyzer {
    return &analyzerService{
        httpClient: httpClient,
        htmlParser: htmlParser,
        logger:     logger,
    }
}

func (s *analyzerService) Analyze(ctx context.Context, urlStr string) (*domain.PageAnalysis, error) {
    // Validate URL
    if _, err := url.ParseRequestURI(urlStr); err != nil {
        s.logger.Error("invalid URL", "url", urlStr, "error", err)
        return nil, domain.ErrInvalidURL
    }

    // Fetch page content
    content, err := s.httpClient.FetchPage(ctx, urlStr)
    if err != nil {
        s.logger.Error("failed to fetch page", "url", urlStr, "error", err)
        switch {
        case err == context.DeadlineExceeded:
            return nil, domain.ErrTimeout
        default:
            return nil, domain.ErrPageNotAccessible
        }
    }


    // Analyze page
    analysis := &domain.PageAnalysis{
        HTMLVersion:  s.htmlParser.GetHTMLVersion(content),
        PageTitle:    s.htmlParser.GetTitle(content),
        Headings:     s.htmlParser.CountHeadings(content),
        Links:        s.htmlParser.AnalyzeLinks(content, urlStr),
        HasLoginForm: s.htmlParser.HasLoginForm(content),
    }

    s.logger.Info("page analysis completed", "url", urlStr)
    return analysis, nil
}