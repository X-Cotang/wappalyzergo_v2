package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	FormatText  OutputFormat = "text"
	FormatJSON  OutputFormat = "json"
	FormatJSONL OutputFormat = "jsonl"
)

// ScanResult represents a single scan result
type ScanResult struct {
	URL          string            `json:"url"`
	Technologies map[string]string `json:"technologies"`
	Mode         string            `json:"mode,omitempty"`
	Error        string            `json:"error,omitempty"`
}

// OutputWriter interface for different output formats
type OutputWriter interface {
	WriteResult(result *ScanResult) error
	Close() error
}

// NewOutputWriter creates an output writer based on format
func NewOutputWriter(format OutputFormat, outputFile string) (OutputWriter, error) {
	var writer io.Writer = os.Stdout
	var closeFunc func() error

	// Open output file if specified
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create output file: %w", err)
		}
		writer = file
		closeFunc = file.Close
	}

	switch format {
	case FormatText:
		return &TextWriter{writer: writer, closeFunc: closeFunc}, nil
	case FormatJSON:
		return &JSONWriter{writer: writer, closeFunc: closeFunc, results: []ScanResult{}}, nil
	case FormatJSONL:
		return &JSONLWriter{writer: writer, closeFunc: closeFunc}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}

// TextWriter writes human-readable text output
type TextWriter struct {
	writer    io.Writer
	closeFunc func() error
}

func (w *TextWriter) WriteResult(result *ScanResult) error {
	if result.Error != "" {
		fmt.Fprintf(w.writer, "\n%s [error]\n  %s\n", result.URL, result.Error)
		return nil
	}

	mode := result.Mode
	if mode == "" {
		mode = "unknown"
	}

	fmt.Fprintf(w.writer, "\n%s [%s]\n", result.URL, mode)
	for tech, version := range result.Technologies {
		if version != "" {
			fmt.Fprintf(w.writer, "  [OK] %s:%s\n", tech, version)
		} else {
			fmt.Fprintf(w.writer, "  [OK] %s\n", tech)
		}
	}

	return nil
}

func (w *TextWriter) Close() error {
	if w.closeFunc != nil {
		return w.closeFunc()
	}
	return nil
}

// JSONWriter collects all results and writes a single JSON array
type JSONWriter struct {
	writer    io.Writer
	closeFunc func() error
	results   []ScanResult
}

func (w *JSONWriter) WriteResult(result *ScanResult) error {
	w.results = append(w.results, *result)
	return nil
}

func (w *JSONWriter) Close() error {
	// Write all results as JSON array
	output := map[string]interface{}{
		"results": w.results,
	}

	encoder := json.NewEncoder(w.writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	if w.closeFunc != nil {
		return w.closeFunc()
	}
	return nil
}

// JSONLWriter writes one JSON object per line
type JSONLWriter struct {
	writer    io.Writer
	closeFunc func() error
}

func (w *JSONLWriter) WriteResult(result *ScanResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	_, err = fmt.Fprintf(w.writer, "%s\n", data)
	return err
}

func (w *JSONLWriter) Close() error {
	if w.closeFunc != nil {
		return w.closeFunc()
	}
	return nil
}

// parseOutputFormat converts string to OutputFormat
func parseOutputFormat(format string) (OutputFormat, error) {
	format = strings.ToLower(strings.TrimSpace(format))

	switch format {
	case "text", "txt", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "jsonl", "ndjson":
		return FormatJSONL, nil
	default:
		return "", fmt.Errorf("unsupported format: %s (supported: text, json, jsonl)", format)
	}
}
