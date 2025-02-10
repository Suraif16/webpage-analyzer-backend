package http

import (
    "context"
    "io"
    "net/http"
    "time"
    "github.com/suraif16/webpage-analyzer/internal/core/domain"
)

type client struct {
    httpClient *http.Client
}

func NewHTTPClient(timeout time.Duration) *client {
    return &client{
        httpClient: &http.Client{
            Timeout: timeout,
            CheckRedirect: func(req *http.Request, via []*http.Request) error {
                if len(via) >= 10 {
                    return domain.ErrPageNotAccessible
                }
                return nil
            },
        },
    }
}

func (c *client) FetchPage(ctx context.Context, url string) (string, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return "", domain.ErrInvalidURL
    }

    // Set user agent to avoid being blocked
    req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WebAnalyzer/1.0)")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return "", domain.ErrPageNotAccessible
    }
    defer resp.Body.Close()

    // Check status code
    switch resp.StatusCode {
    case http.StatusOK:
    case http.StatusNotFound:
        return "", domain.ErrPageNotFound
    default:
        return "", &domain.APIError{
            StatusCode:   resp.StatusCode,
            Message:     resp.Status,
            Description: "Failed to fetch the page",
        }
    }

    // Read body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", domain.ErrInternalServer
    }

    return string(body), nil
}

func (c *client) CheckLink(ctx context.Context, url string) bool {
    req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
    if err != nil {
        return false
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    return resp.StatusCode == http.StatusOK
}