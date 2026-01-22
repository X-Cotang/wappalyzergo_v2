package main

import (
	httputil "github.com/projectdiscovery/wappalyzergo/internal/http"
)

// fetchURLStatic fetches a URL using only HTTP (no JavaScript execution).
// This is faster but cannot detect technologies that rely on JavaScript.
func fetchURLStatic(url string) (map[string][]string, []byte, error) {
	client := httputil.NewClient(*timeout, *userAgent)
	return client.Fetch(url)
}
