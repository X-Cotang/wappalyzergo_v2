package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	wappalyzer "github.com/projectdiscovery/wappalyzergo"
	browserutil "github.com/projectdiscovery/wappalyzergo/internal/browser"
)

// Result holds the scan result for a single URL (simple format)
type Result struct {
	URL          string            `json:"url"`
	Technologies map[string]string `json:"technologies"`
	Mode         string            `json:"mode,omitempty"`
	Error        string            `json:"error,omitempty"`
}

// DetailedResult holds the scan result for a single URL (detailed format)
type DetailedResult struct {
	URL          string                       `json:"url"`
	Technologies map[string]TechnologyDetails `json:"technologies"`
	Mode         string                       `json:"mode,omitempty"`
	Error        string                       `json:"error,omitempty"`
}

// TechnologyDetails contains detailed information about a detected technology
type TechnologyDetails struct {
	Version     string   `json:"version,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	Description string   `json:"description,omitempty"`
	Website     string   `json:"website,omitempty"`
}

func main() {
	// Set up usage function
	flag.Usage = printUsage

	// Parse command-line flags
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("Wappalyzer CLI v%s\n", version)
		return
	}

	// Get URLs from appropriate source (stdin > file > args)
	urls, err := getURLs(*listFile, flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	// Initialize wappalyzer
	var wappalyzerClient *wappalyzer.Wappalyze

	if *fingerprintsDir != "" {
		// Load from directory at runtime (no rebuild needed for changes)
		if !*silent {
			fmt.Fprintf(os.Stderr, "[INFO] Loading fingerprints from: %s\n", *fingerprintsDir)
		}
		wappalyzerClient, err = loadFingerprintsFromDirectory(*fingerprintsDir)
	} else {
		// Use embedded fingerprints
		wappalyzerClient, err = wappalyzer.New()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing wappalyzer: %v\n", err)
		os.Exit(1)
	}

	// Handle legacy -json flag
	outputFormat := *format
	if *jsonOutput {
		outputFormat = "json"
		if !*silent {
			fmt.Fprintln(os.Stderr, "[WARN] -json flag is deprecated, use -format json instead")
		}
	}

	// Parse output format
	parsedFormat, err := parseOutputFormat(outputFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create output writer
	writer, err := NewOutputWriter(parsedFormat, *outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output writer: %v\n", err)
		os.Exit(1)
	}
	defer writer.Close()

	// Scan URLs
	if *detailed {
		processURLsDetailed(urls, wappalyzerClient)
	} else {
		// Use concurrent scanner
		if err := scanURLsConcurrent(urls, wappalyzerClient, writer, *concurrency); err != nil {
			fmt.Fprintf(os.Stderr, "Error during scanning: %v\n", err)
			os.Exit(1)
		}
	}
}

// processURLsDetailed processes URLs and outputs detailed results
func processURLsDetailed(urls []string, wappalyzerClient *wappalyzer.Wappalyze) {
	var results []DetailedResult

	for _, url := range urls {
		mode := "hybrid"
		if *staticMode {
			mode = "static"
		}

		result := DetailedResult{URL: url, Mode: mode}

		// Always use static fetch first
		headers, body, err := fetchURLStatic(url)
		if err != nil {
			result.Error = err.Error()
			if *jsonOutput {
				results = append(results, result)
			} else {
				fmt.Fprintf(os.Stderr, "[!] %s: %v\n", url, err)
			}
			continue
		}

		fingerprints := wappalyzerClient.FingerprintWithInfo(headers, body)
		result.Technologies = formatDetailedFingerprints(fingerprints)

		// Enhance with browser detection if not in static mode
		if !*staticMode {
			// Convert detailed results to simple map for browser enhancement
			simpleTech := make(map[string]string)
			for name, details := range result.Technologies {
				simpleTech[name] = details.Version
			}

			// Setup browser detector
			detector := browserutil.NewDetector(*headless, *userAgent, *waitTime)
			ctx, cancel := context.WithTimeout(context.Background(), *timeout)
			err := detector.EnhanceWithVersions(ctx, url, simpleTech, wappalyzerClient.GetFingerprints())
			cancel()
			if err != nil && !*silent {
				fmt.Fprintf(os.Stderr, "[WARN] Browser detection failed for %s: %v\n", url, err)
			}

			// Update result with enhanced versions
			for name, version := range simpleTech {
				if details, exists := result.Technologies[name]; exists {
					details.Version = version
					result.Technologies[name] = details
				}
			}
		}

		if *jsonOutput {
			results = append(results, result)
		} else {
			printDetailedResult(url, result.Technologies, mode)
		}
	}

	if *jsonOutput {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
	}
}
