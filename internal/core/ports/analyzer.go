package ports

import (
	"context"
	"github.com/suraif16/webpage-analyzer/internal/core/domain"
)

// PageAnalyzer defines the interface for webpage analysis
type PageAnalyzer interface {
	Analyze(ctx context.Context, url string) (*domain.PageAnalysis, error)
}

// HTMLParser defines the interface for HTML parsing operations
type HTMLParser interface {
	GetHTMLVersion(doc string) string
	GetTitle(doc string) string
	CountHeadings(doc string) domain.HeadingCount
	AnalyzeLinks(doc string, baseURL string) domain.LinkAnalysis
	HasLoginForm(doc string) bool
}

// HTTPClient defines the interface for making HTTP requests
type HTTPClient interface {
	FetchPage(ctx context.Context, url string) (string, error)
	CheckLink(ctx context.Context, url string) bool
}
