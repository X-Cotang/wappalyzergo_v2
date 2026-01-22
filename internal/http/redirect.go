package http

import (
	"net/url"
	"strings"
)

// IsSafeDomainRedirect checks if a redirect is to the same domain or www subdomain.
// This prevents following redirects to potentially malicious domains.
//
// Examples of allowed redirects:
//   - example.com -> example.com/path
//   - example.com -> www.example.com
//   - www.example.com -> example.com
//
// Examples of blocked redirects:
//   - example.com -> otherdomain.com
//   - subdomain.example.com -> different.subdomain.example.com
func IsSafeDomainRedirect(original, new *url.URL) bool {
	originalHost := strings.ToLower(original.Host)
	newHost := strings.ToLower(new.Host)

	// Same host is always OK
	if originalHost == newHost {
		return true
	}

	// Remove www. prefix for comparison
	originalBase := strings.TrimPrefix(originalHost, "www.")
	newBase := strings.TrimPrefix(newHost, "www.")

	// Allow if base domains match (example.com <-> www.example.com)
	if originalBase == newBase {
		return true
	}

	// Allow if one has www. and other doesn't but base domain matches
	if originalBase == newHost || newBase == originalHost {
		return true
	}

	return false
}
