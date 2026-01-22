package main

import (
	"fmt"
	"strings"

	wappalyzer "github.com/projectdiscovery/wappalyzergo"
)

// formatSimpleFingerprints converts fingerprint map to simple tech->version map
func formatSimpleFingerprints(fingerprints map[string]struct{}) map[string]string {
	result := make(map[string]string)
	for tech := range fingerprints {
		parts := strings.Split(tech, ":")
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		} else {
			result[tech] = ""
		}
	}
	return result
}

// formatDetailedFingerprints converts fingerprint map with info to detailed result
func formatDetailedFingerprints(fingerprints map[string]wappalyzer.AppInfo) map[string]TechnologyDetails {
	result := make(map[string]TechnologyDetails)
	for tech, info := range fingerprints {
		parts := strings.Split(tech, ":")
		name := parts[0]
		version := ""
		if len(parts) == 2 {
			version = parts[1]
		}

		result[name] = TechnologyDetails{
			Version:     version,
			Categories:  info.Categories,
			Description: info.Description,
			Website:     info.Website,
		}
	}
	return result
}

// printSimpleResult prints technologies in simple format
func printSimpleResult(url string, technologies map[string]string, mode string) {
	fmt.Printf("\n%s [%s]\n", url, mode)
	if len(technologies) == 0 {
		fmt.Println("  No technologies detected")
		return
	}

	for tech, version := range technologies {
		if version != "" {
			fmt.Printf("  [OK] %s:%s\n", tech, version)
		} else {
			fmt.Printf("  [OK] %s\n", tech)
		}
	}
}

// printDetailedResult prints technologies with detailed information
func printDetailedResult(url string, technologies map[string]TechnologyDetails, mode string) {
	fmt.Printf("\n%s [%s]\n", url, mode)
	fmt.Println(strings.Repeat("=", len(url)+len(mode)+3))

	if len(technologies) == 0 {
		fmt.Println("  No technologies detected")
		return
	}

	for tech, details := range technologies {
		fmt.Printf("\n  %s", tech)
		if details.Version != "" {
			fmt.Printf(" (v%s)", details.Version)
		}
		fmt.Println()

		if len(details.Categories) > 0 {
			fmt.Printf("    Categories: %s\n", strings.Join(details.Categories, ", "))
		}

		if details.Website != "" {
			fmt.Printf("    Website: %s\n", details.Website)
		}

		if details.Description != "" {
			// Truncate long descriptions
			desc := details.Description
			if len(desc) > 100 {
				desc = desc[:97] + "..."
			}
			fmt.Printf("    Description: %s\n", desc)
		}
	}
}
