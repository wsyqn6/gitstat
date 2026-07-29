package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	githubAPI   = "https://api.github.com/repos/wsyqn6/gitstat/releases/latest"
	cacheTTL    = 5 * time.Minute
)

type UpdateCheckResponse struct {
	HasUpdate      bool   `json:"hasUpdate"`
	LatestVersion  string `json:"latestVersion"`
	DownloadURL    string `json:"downloadUrl"`
	CurrentVersion string `json:"currentVersion"`
}

type githubRelease struct {
	TagName    string `json:"tag_name"`
	HTMLURL    string `json:"html_url"`
	Prerelease bool   `json:"prerelease"`
}

type cachedResult struct {
	data      UpdateCheckResponse
	expiresAt time.Time
}

var (
	httpClient  = &http.Client{Timeout: 10 * time.Second}
	cacheMu    sync.Mutex
	updateCache *cachedResult
)

func parseVersion(v string) (major, minor, patch int, ok bool) {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return 0, 0, 0, false
	}
	var err error
	major, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, false
	}
	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, false
	}
	patch, err = strconv.Atoi(strings.SplitN(parts[2], "-", 2)[0])
	if err != nil {
		return 0, 0, 0, false
	}
	return major, minor, patch, true
}

func versionGreaterThan(a, b string) bool {
	ma, mi, pa, oka := parseVersion(a)
	mb, mi2, pb, okb := parseVersion(b)
	if !oka || !okb {
		return false
	}
	if ma != mb {
		return ma > mb
	}
	if mi != mi2 {
		return mi > mi2
	}
	return pa > pb
}

func fetchLatestRelease() (*githubRelease, error) {
	req, err := http.NewRequest("GET", githubAPI, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gitstat")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &release, nil
}

func GetUpdateCheckHandler(currentVersion string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cacheMu.Lock()
		cached := updateCache
		cacheMu.Unlock()

		if cached != nil && time.Now().Before(cached.expiresAt) {
			writeJSON(w, "update-check", cached.data)
			return
		}

		release, err := fetchLatestRelease()
		if err != nil {
			writeJSON(w, "update-check", UpdateCheckResponse{
				HasUpdate:      false,
				LatestVersion:  "",
				DownloadURL:    "",
				CurrentVersion: currentVersion,
			})
			return
		}

		latestTag := release.TagName
		hasUpdate := versionGreaterThan(latestTag, currentVersion)

		resp := UpdateCheckResponse{
			HasUpdate:      hasUpdate,
			LatestVersion:  latestTag,
			DownloadURL:    release.HTMLURL,
			CurrentVersion: currentVersion,
		}

		cacheMu.Lock()
		updateCache = &cachedResult{data: resp, expiresAt: time.Now().Add(cacheTTL)}
		cacheMu.Unlock()

		writeJSON(w, "update-check", resp)
	}
}
