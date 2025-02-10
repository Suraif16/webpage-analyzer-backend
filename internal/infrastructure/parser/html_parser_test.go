package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suraif16/webpage-analyzer/internal/core/domain"
)

func TestHTMLParser_GetHTMLVersion(t *testing.T) {
    tests := []struct {
        name     string
        html     string
        expected string
    }{
        {
            name: "HTML5 simple DOCTYPE",
            html: "<!DOCTYPE html><html><head></head><body></body></html>",
            expected: "HTML5",
        },
        {
            name: "HTML5 lowercase DOCTYPE",
            html: "<!doctype html><html></html>",
            expected: "HTML5",
        },
        {
            name: "HTML 4.01 Strict",
            html: `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN">
                  <html></html>`,
            expected: "HTML 4.01",
        },
        {
            name: "HTML 4.01 Transitional",
            html: `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
                  <html></html>`,
            expected: "HTML 4.01",
        },
        {
            name: "XHTML 1.0",
            html: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN">
                  <html></html>`,
            expected: "XHTML",
        },
        {
            name: "No DOCTYPE",
            html: "<html><head></head><body></body></html>",
            expected: "Unknown",
        },
        {
            name: "Malformed DOCTYPE",
            html: "<!DOCTYPE something><html></html>",
            expected: "Unknown",
        },
    }

    parser := NewHTMLParser()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := parser.GetHTMLVersion(tt.html)
            assert.Equal(t, tt.expected, result, 
                "For HTML: %s\nExpected: %s\nGot: %s", 
                tt.html, tt.expected, result)
        })
    }
}

func TestHTMLParser_CountHeadings(t *testing.T) {
    html := `
        <html>
            <body>
                <h1>Title</h1>
                <h2>Subtitle</h2>
                <h2>Another h2</h2>
                <h3>H3 heading</h3>
                <h6>H6 heading</h6>
            </body>
        </html>
    `

    expected := domain.HeadingCount{
        H1: 1,
        H2: 2,
        H3: 1,
        H4: 0,
        H5: 0,
        H6: 1,
    }

    parser := NewHTMLParser()
    result := parser.CountHeadings(html)
    assert.Equal(t, expected, result)
}

func TestHTMLParser_HasLoginForm(t *testing.T) {
    tests := []struct {
        name     string
        html     string
        expected bool
    }{
        {
            name: "Has login form",
            html: `
                <form>
                    <input type="text" name="username">
                    <input type="password" name="password">
                </form>
            `,
            expected: true,
        },
        {
            name: "No login form",
            html: `
                <form>
                    <input type="text" name="search">
                </form>
            `,
            expected: false,
        },
    }

    parser := NewHTMLParser()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := parser.HasLoginForm(tt.html)
            assert.Equal(t, tt.expected, result)
        })
    }
}