package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	wappalyzer "github.com/projectdiscovery/wappalyzergo"
	browserutil "github.com/projectdiscovery/wappalyzergo/internal/browser"
)

func scanURL(url string, wappalyzerClient *wappalyzer.Wappalyze) *ScanResult {
	result := &ScanResult{
		URL:          url,
		Technologies: make(map[string]string),
	}

	// Always use static fetch for initial HTML/headers (fast)
	headers, body, err := fetchURLStatic(url)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// Get technologies from static analysis
	fingerprints := wappalyzerClient.Fingerprint(headers, body)
	result.Technologies = formatSimpleFingerprints(fingerprints)

	// Enhance with browser-based detection if not in static mode
	if !*staticMode {
		// Setup browser detector
		detector := browserutil.NewDetector(*headless, *userAgent, *waitTime)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()

		// Enhance with browser detection (batch execution)
		err := detector.EnhanceWithVersions(ctx, url, result.Technologies, wappalyzerClient.GetFingerprints())
		if err != nil {
			// Browser detection failed, but keep static results
			if !*silent {
				fmt.Fprintf(os.Stderr, "[WARN] Browser detection failed for %s: %v\n", url, err)
			}
		}
		result.Mode = "hybrid"
	} else {
		result.Mode = "static"
	}

	return result
}

// scanURLsConcurrent scans multiple URLs with concurrency control
func scanURLsConcurrent(urls []string, wappalyzerClient *wappalyzer.Wappalyze, writer OutputWriter, concurrency int) error {
	// Create channel for URLs to scan
	urlChan := make(chan string, len(urls))
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	// Create wait group and result channel
	var wg sync.WaitGroup
	resultChan := make(chan *ScanResult, concurrency)
	errorChan := make(chan error, 1)

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urlChan {
				result := scanURL(url, wappalyzerClient)
				resultChan <- result
			}
		}()
	}

	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect and write results
	count := 0
	for result := range resultChan {
		if err := writer.WriteResult(result); err != nil {
			select {
			case errorChan <- fmt.Errorf("failed to write result: %w", err):
			default:
			}
			break
		}
		count++

		// Progress indicator (only if not silent and outputting to file)
		if !*silent && *outputFile != "" {
			fmt.Fprintf(os.Stderr, "\r[INFO] Processed %d/%d URLs", count, len(urls))
		}
	}

	if !*silent && *outputFile != "" && count > 0 {
		fmt.Fprintln(os.Stderr) // New line after progress
	}

	// Check for errors
	select {
	case err := <-errorChan:
		return err
	default:
		return nil
	}
}
