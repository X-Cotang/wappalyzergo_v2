package wappalyzer

import (
	"embed"
	"encoding/json"
	"io/fs"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	//go:embed data
	dataFS embed.FS

	syncOnce          sync.Once
	categoriesMapping map[int]categoryItem
	fingerprintsCache string
	fingerprints      string // Deprecated: Use fingerprintsCache or GetRawFingerprints()
)

func init() {
	syncOnce.Do(func() {
		categoriesMapping = make(map[int]categoryItem)

		// Load categories
		categoriesData, err := dataFS.ReadFile("data/categories/categories.json")
		if err != nil {
			panic("failed to load categories: " + err.Error())
		}

		var categories map[string]categoryItem
		if err := json.Unmarshal(categoriesData, &categories); err != nil {
			panic("failed to parse categories: " + err.Error())
		}

		for category, data := range categories {
			parsed, _ := strconv.Atoi(category)
			categoriesMapping[parsed] = data
		}

		// Load and merge all fingerprint files
		fingerprintsCache = loadAllFingerprints()
		fingerprints = fingerprintsCache // For backward compatibility
	})
}

func loadAllFingerprints() string {
	allApps := make(map[string]interface{})

	// Walk through all JSON files in data/fingerprints/
	err := fs.WalkDir(dataFS, "data/fingerprints", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-JSON files
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		// Read the category file
		data, err := dataFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Parse the category file
		var categoryData struct {
			Apps map[string]interface{} `json:"apps"`
		}
		if err := json.Unmarshal(data, &categoryData); err != nil {
			return err
		}

		// Merge apps into allApps
		for appName, appData := range categoryData.Apps {
			allApps[appName] = appData
		}

		return nil
	})

	if err != nil {
		panic("failed to load fingerprints: " + err.Error())
	}

	// Create the expected format
	result := map[string]interface{}{
		"apps": allApps,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		panic("failed to marshal fingerprints: " + err.Error())
	}

	return string(resultJSON)
}

func GetRawFingerprints() string {
	return fingerprintsCache
}

func GetCategoriesMapping() map[int]categoryItem {
	return categoriesMapping
}

type categoryItem struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}
