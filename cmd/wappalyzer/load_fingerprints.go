package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	wappalyzer "github.com/projectdiscovery/wappalyzergo"
)

// loadFingerprintsFromDirectory loads all JSON files from a directory
func loadFingerprintsFromDirectory(dirPath string) (*wappalyzer.Wappalyze, error) {
	files, err := filepath.Glob(filepath.Join(dirPath, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list fingerprint files: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no fingerprint files found in %s", dirPath)
	}

	// Load and merge all JSON files
	mergedApps := make(map[string]*wappalyzer.Fingerprint)

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", file, err)
		}

		var fps wappalyzer.Fingerprints
		if err := json.Unmarshal(data, &fps); err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", file, err)
		}

		// Merge apps
		for name, fp := range fps.Apps {
			mergedApps[name] = fp
		}
	}

	// Create final fingerprints structure
	finalFps := &wappalyzer.Fingerprints{
		Apps: mergedApps,
	}

	// Create temporary merged file
	tempFile := filepath.Join(os.TempDir(), "merged_fingerprints.json")
	tempData, err := json.Marshal(finalFps)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merged fingerprints: %w", err)
	}

	if err := os.WriteFile(tempFile, tempData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	defer os.Remove(tempFile)

	// Use NewFromFile with the merged data
	return wappalyzer.NewFromFile(tempFile, false, false)
}
