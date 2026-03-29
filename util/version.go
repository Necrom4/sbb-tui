package utils

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

func LatestVersion() (string, error) {
	url := "https://api.github.com/repos/Necrom4/sbb-tui/releases/latest"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "sbb-tui")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api returned %s", resp.Status)
	}

	var release releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func NewerVersion(current string) (string, error) {
	if current == "dev" {
		return "", nil
	}

	latest, err := LatestVersion()
	if err != nil {
		return "", err
	}

	latestValid := semver.IsValid(latest)
	if !latestValid {
		return "", fmt.Errorf("latest version is not valid: %s", latest)
	}

	currentValid := semver.IsValid(current)
	if !currentValid {
		return "", fmt.Errorf("current version is not valid: %s", current)
	}

	cmp := semver.Compare(current, latest)

	if cmp < 0 {
		return latest, nil
	}

	return "", nil
}
