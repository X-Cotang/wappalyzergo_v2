package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// readURLsFromStdin reads URLs from stdin (pipe input)
func readURLsFromStdin() ([]string, error) {
	var urls []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from stdin: %w", err)
	}

	return urls, nil
}

// readURLsFromFile reads URLs from a file (one per line)
func readURLsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return urls, nil
}

// hasStdinData checks if stdin has data available (piped input)
func hasStdinData() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// getURLs returns URLs from the appropriate source based on flags and stdin
// Priority: stdin > file list > command-line args
func getURLs(listFile string, args []string) ([]string, error) {
	// Check for stdin input first (highest priority)
	if hasStdinData() {
		if !*silent {
			fmt.Fprintln(os.Stderr, "[INFO] Reading URLs from stdin...")
		}
		return readURLsFromStdin()
	}

	// Check for file list
	if listFile != "" {
		if !*silent {
			fmt.Fprintf(os.Stderr, "[INFO] Reading URLs from file: %s\n", listFile)
		}
		return readURLsFromFile(listFile)
	}

	// Use command-line arguments
	if len(args) > 0 {
		return args, nil
	}

	// No input source
	return nil, fmt.Errorf("no URLs provided (use command-line args, -l flag, or pipe input)")
}
