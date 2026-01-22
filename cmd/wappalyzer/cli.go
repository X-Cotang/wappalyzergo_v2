package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// CLI flags - all flag variables are defined here for easy management
var (
	// Output flags
	jsonOutput = flag.Bool("json", false, "Output results as JSON (deprecated, use -format json)")
	detailed   = flag.Bool("detailed", false, "Show detailed information including categories and descriptions")
	format     = flag.String("format", "text", "Output format: text, json, jsonl")
	outputFile = flag.String("o", "", "Write output to file instead of stdout")
	silent     = flag.Bool("silent", false, "Silent mode, only output results")

	// Detection mode flags
	staticMode = flag.Bool("static", false, "Use static HTTP mode (no JavaScript execution)")
	headless   = flag.Bool("headless", true, "Run browser in headless mode")

	// Timing flags
	waitTime = flag.Duration("wait", 3*time.Second, "Wait time for JavaScript execution in browser mode")
	timeout  = flag.Duration("timeout", 30*time.Second, "Total timeout for fetching URL")

	// Input flags
	listFile = flag.String("l", "", "Read URLs from file (one per line)")

	// Performance flags
	concurrency = flag.Int("c", 1, "Number of concurrent requests")

	// Configuration flags
	userAgent       = flag.String("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0", "Custom User-Agent header")
	fingerprintsDir = flag.String("fingerprints-dir", "", "Load fingerprints from directory instead of embedded data (e.g., ./data/fingerprints)")

	// Info flags
	showVersion = flag.Bool("version", false, "Show version information")
)

const version = "2.0.0"

// printUsage prints the command usage information
func printUsage() {
	fmt.Fprintf(os.Stderr, "Wappalyzer CLI v%s - Detect web technologies\n\n", version)
	fmt.Fprintf(os.Stderr, "Usage: %s [options] <url> [url2] [url3]...\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "By default, uses headless browser for accurate JavaScript-based detection.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s https://nextjs.org/\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -json https://nextjs.org/\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -detailed https://nextjs.org/\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -static https://nextjs.org/  # Fast mode, no JS\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -c 5 -l urls.txt  # Concurrent scanning\n", os.Args[0])
}
