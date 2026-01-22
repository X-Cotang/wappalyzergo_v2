package browser

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	wappalyzer "github.com/projectdiscovery/wappalyzergo"
)

// Detector handles browser-based technology detection and version extraction.
type Detector struct {
	headless  bool
	userAgent string
	waitTime  time.Duration
}

// NewDetector creates a new browser detector with the specified configuration.
func NewDetector(headless bool, userAgent string, waitTime time.Duration) *Detector {
	return &Detector{
		headless:  headless,
		userAgent: userAgent,
		waitTime:  waitTime,
	}
}

func (d *Detector) EnhanceWithVersions(
	ctx context.Context,
	url string,
	technologies map[string]string,
	fingerprints *wappalyzer.Fingerprints,
) error {
	// Setup browser context
	browserCtx, cancel := SetupContext(ctx, d.headless, d.userAgent)
	defer cancel()

	// Navigate and wait for page to be ready
	err := chromedp.Run(browserCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Sleep(d.waitTime),
	)

	if err != nil {
		return err
	}

	// Execute detection and version extraction for each app individually
	// This avoids JavaScript syntax errors from special characters in variable names
	for appName, fingerprint := range fingerprints.Apps {
		// Skip if no browser detection configured

		if fingerprint.Browser == nil {
			continue
		}

		// Check if already detected by static analysis
		_, alreadyDetected := technologies[appName]

		// Run detection if not already detected
		detected := alreadyDetected
		if !detected && len(fingerprint.Browser.Detection) > 0 {
			detected = ExecuteDetectionRules(browserCtx, fingerprint.Browser.Detection)
		}

		// If detected (either already or via browser), try to extract version
		if detected {
			version := ""

			// If version extraction rules exist, try them
			if len(fingerprint.Browser.Version) > 0 {
				version = ExtractVersion(browserCtx, fingerprint.Browser.Version)
			}

			// Update or add technology
			if version != "" {
				technologies[appName] = version
			} else if !alreadyDetected {
				// Add without version if detected but no version found
				technologies[appName] = ""
			}
		}
	}

	return nil
}
