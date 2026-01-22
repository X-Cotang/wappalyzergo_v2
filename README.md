# Wappalyzergo

A high performance port of the Wappalyzer Technology Detection Library to Go. Inspired by [Webanalyze](https://github.com/rverton/webanalyze).

Uses data from 
- https://github.com/enthec/webappanalyzer
- https://github.com/HTTPArchive/wappalyzer

## Features

### Core Library
- Very simple and easy to use, with clean codebase.
- Normalized regexes + auto-updating database of wappalyzer fingerprints.
- Optimized for performance: parsing HTML manually for best speed.
- Well-organized fingerprint data in `data/fingerprints/` directory for easy management.

### CLI Tool
- 🚀 **Powerful command-line interface** with multiple output formats (text, JSON, JSONL)
- 🌐 **Hybrid detection mode**: Combines fast HTTP analysis with browser-based JavaScript execution
- 🔍 **Accurate version detection**: Uses chromedp to execute JavaScript and extract framework versions
- ⚡ **Concurrent scanning**: Process multiple URLs in parallel with `-c` flag
- 📝 **Flexible input**: Read from stdin, files, or command-line arguments
- 🎯 **Smart redirect handling**: Safely follows same-domain redirects

### Browser Integration (chromedp)
- Execute JavaScript for accurate detection of client-side technologies
- Extract precise version numbers from frameworks (Next.js, React, Vue, Angular)
- Support for both static mode (fast) and hybrid mode (accurate)
- Detect technologies that only exist in browser runtime

### Using *go install*

```sh
go install -v github.com/projectdiscovery/wappalyzergo/cmd/update-fingerprints@latest
```

After this command *wappalyzergo* library source will be in your current go.mod.

## Example
Usage Example:

``` go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	wappalyzer "github.com/projectdiscovery/wappalyzergo"
)

func main() {
	resp, err := http.DefaultClient.Get("https://www.hackerone.com")
	if err != nil {
		log.Fatal(err)
	}
	data, _ := io.ReadAll(resp.Body) // Ignoring error for example

	wappalyzerClient, err := wappalyzer.New()
	fingerprints := wappalyzerClient.Fingerprint(resp.Header, data)
	fmt.Printf("%v\n", fingerprints)

	// Output: map[Acquia Cloud Platform:{} Amazon EC2:{} Apache:{} Cloudflare:{} Drupal:{} PHP:{} Percona:{} React:{} Varnish:{}]
}
```

## Command Line Tool

### Installation

```sh
go install github.com/projectdiscovery/wappalyzergo/cmd/wappalyzer@latest
```

Or build from source:

```sh
git clone https://github.com/projectdiscovery/wappalyzergo.git
cd wappalyzergo
go build -o wappalyzer ./cmd/wappalyzer
```

### Usage

**Basic usage:**
```sh
$ wappalyzer https://nextjs.org/

https://nextjs.org/
  [OK] Next.js
  [OK] React
  [OK] Webpack
  [OK] Vercel
  [OK] HSTS
  [OK] Node.js
```

**Detailed output with categories and descriptions:**
```sh
$ wappalyzer -detailed https://nextjs.org/

https://nextjs.org/
===================

  Next.js
    Categories: JavaScript frameworks, Web frameworks
    Website: https://nextjs.org
    Description: Next.js is a React framework for developing single page Javascript applications.

  React
    Categories: JavaScript frameworks
    Website: https://reactjs.org
    Description: React is an open-source JavaScript library for building user interfaces or UI components.
  
  ...
```

**JSON output:**
```sh
$ wappalyzer -json https://nextjs.org/
[
  {
    "url": "https://nextjs.org/",
    "technologies": {
      "Next.js": "",
      "React": "",
      "Webpack": "",
      "Vercel": "",
      "HSTS": "",
      "Node.js": ""
    }
  }
]
```

**Multiple URLs:**
```sh
$ wappalyzer https://nextjs.org/ https://react.dev/
```

### CLI Options

| Flag | Description | Default |
|------|-------------|---------|
| `-json` | Output results as JSON | `false` |
| `-detailed` | Show detailed information | `false` |
| `-static` | Use static-only mode (no browser) | `false` |
| `-wait` | JavaScript execution wait time | `3s` |
| `-timeout` | Total request timeout | `30s` |
| `-user-agent` | Custom User-Agent header | `Mozilla/5.0...` |
| `-headless` | Run browser in headless mode | `true` |
| `-version` | Show version information | - |

### Detection Modes

**Hybrid Mode (Default)**:
```sh
wappalyzer https://nextjs.org/
# Output: Next.js:16.1.1 (with version detected via browser)
```
- Fast static HTTP analysis for comprehensive technology detection
- Browser-based JavaScript execution ONLY for version extraction
- Best of both worlds: speed + accuracy

**Static-Only Mode**:
```sh
wappalyzer -static https://nextjs.org/
# Output: Next.js (without version)
```
- Fastest mode, pure HTTP analysis
- Use when Chrome/Chromium not available
- Still detects most technologies

**How Hybrid Mode Works**:
1. Fetch page via HTTP (fast)
2. Static analysis with wappalyzergo
3. If framework detected without version → launch browser
4. Extract JavaScript variables (next.version, React.version, etc.)
5. Merge results

**Supported Frameworks for Version Detection**:
- Next.js (`next.version`)
- React (`React.version`)
- Vue.js (`Vue.version`)
- Angular (`angular.version.full`)

### Redirect Handling

The CLI automatically follows redirects with smart security policies:

**✅ Allowed Redirects**:
- Same domain: `example.com` → `example.com/blog`
- www subdomain: `example.com` ↔ `www.example.com`
- Path changes within domain

**❌ Blocked Redirects**:
- Different domains: `test.com` → `otherdomain.com`
- Different subdomains (except www)
- Maximum 10 redirects to prevent loops

**Example**:
```sh
# Follows redirect from github.com to www.github.com
wappalyzer https://github.com
```

## Project Improvements

This fork includes several enhancements over the original wappalyzergo:

### 1. 🛠️ Enhanced CLI Tool

**Added comprehensive command-line interface** with features:
- Multiple output formats: `text`, `json`, `jsonl`
- Concurrent URL scanning with `-c` flag
- Flexible input from stdin, files, or arguments
- File output with `-o` flag
- Silent mode for scripting

**Example - Concurrent scanning**:
```sh
# Scan 5 URLs concurrently, output to file
wappalyzer -c 5 -l urls.txt -o results.json -format jsonl
```

### 2. 📁 Better Fingerprint Organization

**Reorganized fingerprint data** for easier management:
- Split monolithic `fingerprints_data.json` into modular files
- Organized by categories in `data/fingerprints/` directory:
  - `001-cms.json` - Content Management Systems
  - `012-javascript-frameworks.json` - JavaScript frameworks
  - `059-javascript-libraries.json` - JavaScript libraries
  - ...and many more

**Benefits**:
- ✅ Easier to find and update specific fingerprints
- ✅ Better version control with smaller, focused files
- ✅ Reduced merge conflicts when contributing
- ✅ Can load custom fingerprints with `-fingerprints-dir` flag

**Custom fingerprints**:
```sh
# Use custom fingerprint directory
wappalyzer -fingerprints-dir ./my-custom-fingerprints https://example.com
```

### 3. 🌐 Chromedp Integration

**Added headless browser support** for accurate JavaScript-based detection:

**What it enables**:
- 🎯 **Version detection**: Extract exact versions from JavaScript variables
  - `Next.js:16.1.1` instead of just `Next.js`
  - `React:18.2.0` instead of just `React`
- 🔍 **Runtime detection**: Detect technologies that only exist after JavaScript execution
- ✅ **Better accuracy**: Distinguish similar libraries (Lodash vs Underscore.js)

**Architecture**:
- **Hybrid mode** (default): HTTP analysis + selective browser execution
  - Fast static analysis first
  - Browser used ONLY for version extraction
  - Best balance of speed and accuracy

- **Static mode** (`-static`): HTTP-only analysis  
  - Fastest performance
  - No browser dependency
  - Good for bulk scanning

**Internal packages** for reusability:
```go
import (
    "github.com/projectdiscovery/wappalyzergo/internal/browser"
    "github.com/projectdiscovery/wappalyzergo/internal/http"
)

// Use browser detector in your own code
detector := browser.NewDetector(true, userAgent, 3*time.Second)
detector.EnhanceWithVersions(ctx, url, technologies, fingerprints)
```

**Version detection example**:
```sh
# Static mode - no version
$ wappalyzer -static https://nextjs.org/
Next.js

# Hybrid mode - with version  
$ wappalyzer https://nextjs.org/
Next.js:16.1.1  ✓ Accurate version detected!
```

### 4. 🏗️ Improved Code Organization

**Refactored for maintainability**:
- Split large `main.go` (621 lines → 210 lines)
- Created internal packages (`browser`, `http`) for reusability
- Modular CLI structure with focused files:
  - `cli.go` - Flag definitions
  - `browser.go` - Browser detection
  - `static.go` - HTTP fetching  
  - `format.go` - Output formatting

## Quick Start

**Install CLI**:
```sh
go install github.com/projectdiscovery/wappalyzergo/cmd/wappalyzer@latest
```

**Basic scan**:
```sh
wappalyzer https://example.com
```

**Advanced usage**:
```sh
# Scan multiple URLs concurrently with JSON output
cat urls.txt | wappalyzer -c 10 -format jsonl -o results.jsonl

# Static mode for fast bulk scanning
wappalyzer -static -c 20 -l large-url-list.txt

# Detailed output with descriptions
wappalyzer -detailed https://example.com
```

