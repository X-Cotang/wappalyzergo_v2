package browser

import (
	"context"
	"fmt"
	"regexp"

	"github.com/chromedp/chromedp"
	wappalyzer "github.com/projectdiscovery/wappalyzergo"
)

// ExecuteDetectionRules runs browser detection rules and returns true if ANY rule matches.
// Rules can be DOM selectors or JavaScript eval expressions.
func ExecuteDetectionRules(ctx context.Context, rules []wappalyzer.DetectionRule) bool {
	for _, rule := range rules {
		var result bool

		switch rule.Type {
		case "dom-selector":
			// Check if DOM element exists
			err := chromedp.Run(ctx,
				chromedp.Evaluate(fmt.Sprintf(`!!document.querySelector('%s')`, rule.Selector), &result),
			)
			if err == nil && result {
				return true
			}

		case "js-eval":
			// Execute JavaScript and check for truthy result
			err := chromedp.Run(ctx,
				chromedp.Evaluate(rule.Query, &result),
			)
			if err == nil && result {
				return true
			}
		}
	}

	return false
}

// ExtractVersion tries version extraction rules in order and returns the first successful result.
// Supports both DOM attribute extraction and JavaScript evaluation.
func ExtractVersion(ctx context.Context, rules []wappalyzer.VersionExtraction) string {
	for _, rule := range rules {
		var version string

		switch rule.Type {
		case "dom-attribute":
			// Get attribute value from DOM element
			query := fmt.Sprintf(`
				(() => {
					const el = document.querySelector('%s');
					if (el) {
						const value = el.getAttribute('%s');
						if (value) {
							return value;
						}
					}
					return '';
				})()
			`, rule.Selector, rule.Attribute)

			err := chromedp.Run(ctx,
				chromedp.Evaluate(query, &version),
			)

			if err == nil && version != "" {
				// Apply pattern extraction if specified
				if rule.Pattern != "" {
					re, err := regexp.Compile(rule.Pattern)
					if err == nil {
						matches := re.FindStringSubmatch(version)
						if len(matches) > 1 {
							return matches[1] // Return first capture group
						}
					}
				}
				return version
			}

		case "js-eval":
			// Execute JavaScript to get version
			err := chromedp.Run(ctx,
				chromedp.Evaluate(rule.Query, &version),
			)

			if err == nil && version != "" {
				return version
			}
		}
	}

	return ""
}
