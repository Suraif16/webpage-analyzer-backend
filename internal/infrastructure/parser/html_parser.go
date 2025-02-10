package parser

import (
    "strings"
    "net/url"
    "github.com/PuerkitoBio/goquery"
    "github.com/suraif16/webpage-analyzer/internal/core/domain"
)

type htmlParser struct{}

func NewHTMLParser() *htmlParser {
    return &htmlParser{}
}

func (p *htmlParser) GetHTMLVersion(doc string) string {
    docReader := strings.NewReader(doc)
    docParsed, err := goquery.NewDocumentFromReader(docReader)
    if err != nil {
        return "Unknown"
    }

    // Check doctype node
    doctype := docParsed.Find("DOCTYPE").Text()
    switch {
    case strings.Contains(strings.ToLower(doctype), "html 5"):
        return "HTML5"
    case strings.Contains(doc, "<!DOCTYPE html>"):
        return "HTML5"
    case strings.Contains(strings.ToLower(doctype), "html 4.01"):
        return "HTML 4.01"
    case strings.Contains(strings.ToLower(doctype), "xhtml"):
        return "XHTML"
    default:
        return "Unknown"
    }
}

func (p *htmlParser) GetTitle(doc string) string {
    docReader := strings.NewReader(doc)
    docParsed, err := goquery.NewDocumentFromReader(docReader)
    if err != nil {
        return ""
    }
    return docParsed.Find("title").First().Text()
}

func (p *htmlParser) CountHeadings(doc string) domain.HeadingCount {
    docReader := strings.NewReader(doc)
    docParsed, err := goquery.NewDocumentFromReader(docReader)
    if err != nil {
        return domain.HeadingCount{}
    }

    headings := domain.HeadingCount{}
    
    docParsed.Find("h1").Each(func(i int, s *goquery.Selection) {
        headings.H1++
    })
    docParsed.Find("h2").Each(func(i int, s *goquery.Selection) {
        headings.H2++
    })
    docParsed.Find("h3").Each(func(i int, s *goquery.Selection) {
        headings.H3++
    })
    docParsed.Find("h4").Each(func(i int, s *goquery.Selection) {
        headings.H4++
    })
    docParsed.Find("h5").Each(func(i int, s *goquery.Selection) {
        headings.H5++
    })
    docParsed.Find("h6").Each(func(i int, s *goquery.Selection) {
        headings.H6++
    })

    return headings
}

func (p *htmlParser) AnalyzeLinks(doc string, baseURL string) domain.LinkAnalysis {
    docReader := strings.NewReader(doc)
    docParsed, err := goquery.NewDocumentFromReader(docReader)
    if err != nil {
        return domain.LinkAnalysis{}
    }

    baseURLParsed, err := url.Parse(baseURL)
    if err != nil {
        return domain.LinkAnalysis{}
    }

    analysis := domain.LinkAnalysis{}
    
    docParsed.Find("a[href]").Each(func(i int, s *goquery.Selection) {
        href, exists := s.Attr("href")
        if !exists {
            return
        }

        linkURL, err := url.Parse(href)
        if err != nil {
            analysis.Inaccessible++
            return
        }

        // Resolve relative URLs
        linkURL = baseURLParsed.ResolveReference(linkURL)

        if linkURL.Host == baseURLParsed.Host {
            analysis.Internal++
        } else {
            analysis.External++
        }
    })

    return analysis
}

func (p *htmlParser) HasLoginForm(doc string) bool {
    docReader := strings.NewReader(doc)
    docParsed, err := goquery.NewDocumentFromReader(docReader)
    if err != nil {
        return false
    }

    // Check for forms with password fields
    hasLoginForm := false
    docParsed.Find("form").Each(func(i int, s *goquery.Selection) {
        if s.Find("input[type='password']").Length() > 0 {
            hasLoginForm = true
        }
    })

    return hasLoginForm
}