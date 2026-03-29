// Package api provides HTTP client functions for the Swiss public transport API.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/necrom4/sbb-tui/model"
)

type locationsResponse struct {
	Stations []struct {
		Name string `json:"name"`
	} `json:"stations"`
}

type connectionsResponse struct {
	Connections []model.Connection `json:"connections"`
}

// FetchLocations returns station name suggestions matching the given query.
func FetchLocations(query string) ([]string, error) {
	apiURL := "https://transport.opendata.ch/v1/locations?type=station&query=" + url.QueryEscape(query)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result locationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	names := make([]string, 0, len(result.Stations))
	for _, s := range result.Stations {
		if s.Name != "" {
			names = append(names, s.Name)
		}
	}
	return names, nil
}

// FetchConnections queries the transport API for connections between two stations.
func FetchConnections(from, to, date, timeStr string, isArrivalTime bool, limit int) ([]model.Connection, error) {
	parts := []string{
		fmt.Sprintf("from=%s", url.QueryEscape(from)),
		fmt.Sprintf("to=%s", url.QueryEscape(to)),
	}

	if date != "" {
		parts = append(parts, fmt.Sprintf("date=%s", url.QueryEscape(date)))
	}

	if timeStr != "" {
		parts = append(parts, fmt.Sprintf("time=%s", url.QueryEscape(timeStr)))
	}

	isArrival := "0"
	if isArrivalTime {
		isArrival = "1"
	}

	parts = append(parts,
		fmt.Sprintf("isArrivalTime=%s", isArrival),
		fmt.Sprintf("limit=%v", limit),
	)

	apiURL := "https://transport.opendata.ch/v1/connections?" + strings.Join(parts, "&")

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result connectionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Connections, nil
}
