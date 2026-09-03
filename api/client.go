// Package api wraps the search.ch timetable HTTP endpoints.
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/necrom4/sbb-tui/model"
)

// baseURL is a variable so tests can point the client at a test server.
var baseURL = "https://timetable.search.ch/api"

// httpClient bounds every API call so a dead connection (tunnels,
// spotty cellular) fails fast instead of hanging the search forever.
var httpClient = &http.Client{Timeout: 15 * time.Second}

// completionEntry is one suggestion returned by completion.json.
type completionEntry struct {
	Label     string `json:"label"`
	IconClass string `json:"iconclass"`
}

// routeResponse is the top-level shape of route.json. Errors arrive as
// human-readable strings in messages with an HTTP 200 status.
type routeResponse struct {
	Connections []model.Connection `json:"connections"`
	Messages    []string           `json:"messages"`
}

// getJSON fetches apiURL (asking for English texts) and decodes the body into out.
func getJSON(apiURL string, out any) error {
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Language", "en")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned %s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}

// FetchLocations returns station name suggestions matching query.
func FetchLocations(query string) ([]string, error) {
	apiURL := baseURL + "/completion.json?term=" + url.QueryEscape(query)

	var entries []completionEntry
	if err := getJSON(apiURL, &entries); err != nil {
		return nil, fmt.Errorf("fetching locations: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		// Skip plain addresses so suggestions stay stations/stops.
		if e.Label == "" || e.IconClass == "sl-icon-type-adr" {
			continue
		}
		names = append(names, e.Label)
	}
	return names, nil
}

// FetchConnections returns up to `limit` connections between two stations.
// date is expected as YYYY-MM-DD and timeStr as HH:MM.
func FetchConnections(from, to string, via []string, date, timeStr string, isArrivalTime bool, limit int) ([]model.Connection, error) {
	params := url.Values{}
	params.Set("from", from)
	params.Set("to", to)
	for _, v := range via {
		if v != "" {
			params.Add("via[]", v)
		}
	}
	if date != "" {
		params.Set("date", date)
	}
	if timeStr != "" {
		params.Set("time", timeStr)
	}
	if isArrivalTime {
		params.Set("time_type", "arrival")
	}
	params.Set("num", strconv.Itoa(limit))
	params.Set("show_delays", "1")
	params.Set("show_trackchanges", "1")

	apiURL := baseURL + "/route.json?" + params.Encode()

	var result routeResponse
	if err := getJSON(apiURL, &result); err != nil {
		return nil, fmt.Errorf("fetching connections: %w", err)
	}

	if len(result.Connections) == 0 && len(result.Messages) > 0 {
		return nil, errors.New(strings.Join(result.Messages, " "))
	}

	return result.Connections, nil
}
