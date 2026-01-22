package http

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client wraps http.Client with custom redirect handling for safe URL scanning.
type Client struct {
	*http.Client
	userAgent string
}

// NewClient creates a new HTTP client with safe redirect policy.
// The client will only follow redirects to the same domain or www subdomain,
// and will stop after 10 redirects to prevent infinite loops.
func NewClient(timeout time.Duration, userAgent string) *Client {
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Limit redirect count
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}

			// Get original URL (first request)
			originalURL := via[0].URL
			newURL := req.URL

			// Only allow redirects to same domain or www subdomain
			if !IsSafeDomainRedirect(originalURL, newURL) {
				return fmt.Errorf("redirect to different domain not allowed: %s -> %s",
					originalURL.Host, newURL.Host)
			}

			return nil
		},
	}

	return &Client{
		Client:    client,
		userAgent: userAgent,
	}
}

// Fetch fetches a URL and returns the response headers and body.
// The User-Agent header is automatically set to the client's configured value.
func (c *Client) Fetch(url string) (map[string][]string, []byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid URL: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response: %w", err)
	}

	return resp.Header, body, nil
}
