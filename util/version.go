package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/mod/semver"
)

type releaseResponse struct {
	TagName string `json:"tag_name"`
}

var latestReleaseURL = "https://api.github.com/repos/Necrom4/sbb-tui/releases/latest"

func latestVersion() (string, error) {
	req, err := http.NewRequest(http.MethodGet, latestReleaseURL, nil)
	if err != nil {
		return "", fmt.Errorf("checking latest version: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "sbb-tui")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("checking latest version: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("checking latest version: GitHub API returned %s", resp.Status)
	}

	var release releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("checking latest version: decoding response: %w", err)
	}

	return release.TagName, nil
}

func NewerVersion(current string) (string, error) {
	if current == "dev" {
		return "", nil
	}

	latest, err := latestVersion()
	if err != nil {
		return "", fmt.Errorf("checking for newer version: %w", err)
	}

	if !semver.IsValid(latest) {
		return "", fmt.Errorf("checking for newer version: latest tag is not valid semver: %s", latest)
	}

	if !semver.IsValid(current) {
		return "", fmt.Errorf("checking for newer version: current version is not valid semver: %s", current)
	}

	cmp := semver.Compare(current, latest)

	if cmp < 0 {
		return latest, nil
	}

	return "", nil
}
