package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// withTestServer points the package at a test server for the duration of a test.
func withTestServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	orig := baseURL
	baseURL = srv.URL
	t.Cleanup(func() { baseURL = orig })
}

func TestFetchLocations(t *testing.T) {
	var gotPath, gotTerm, gotLang string
	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotTerm = r.URL.Query().Get("term")
		gotLang = r.Header.Get("Accept-Language")
		_, _ = w.Write([]byte(`[
			{"label": "Lausanne", "iconclass": "sl-icon-type-train"},
			{"label": "Lausanne, Riponne", "iconclass": "sl-icon-type-tram"},
			{"label": "Lausannestr. 5", "iconclass": "sl-icon-type-adr"},
			{"label": "", "iconclass": "sl-icon-type-train"}
		]`))
	})

	names, err := FetchLocations("lau")
	if err != nil {
		t.Fatal(err)
	}

	if gotPath != "/completion.json" {
		t.Errorf("path = %s, want /completion.json", gotPath)
	}
	if gotTerm != "lau" {
		t.Errorf("term = %q, want lau", gotTerm)
	}
	if gotLang != "en" {
		t.Errorf("Accept-Language = %q, want en", gotLang)
	}

	want := []string{"Lausanne", "Lausanne, Riponne"}
	if len(names) != len(want) {
		t.Fatalf("names = %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("names = %v, want %v", names, want)
		}
	}
}

func TestFetchLocationsHTTPError(t *testing.T) {
	withTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	if _, err := FetchLocations("lau"); err == nil {
		t.Fatal("expected error on HTTP 500")
	}
}

func TestFetchLocationsBadJSON(t *testing.T) {
	withTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	})

	if _, err := FetchLocations("lau"); err == nil {
		t.Fatal("expected decode error")
	}
}

func TestFetchConnections(t *testing.T) {
	var gotQuery map[string][]string
	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/route.json" {
			t.Errorf("path = %s, want /route.json", r.URL.Path)
		}
		gotQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{"connections": [
			{"from": "Bern", "to": "Luzern", "duration": 3600,
			 "legs": [{"type": "express_train", "line": "IC 6", "name": "Bern"}]}
		]}`))
	})

	conns, err := FetchConnections("Bern", "Luzern", nil, "2026-06-15", "14:00", true, 5)
	if err != nil {
		t.Fatal(err)
	}

	wantParams := map[string]string{
		"from":              "Bern",
		"to":                "Luzern",
		"date":              "2026-06-15",
		"time":              "14:00",
		"time_type":         "arrival",
		"num":               "5",
		"show_delays":       "1",
		"show_trackchanges": "1",
	}
	for k, want := range wantParams {
		if got := gotQuery[k]; len(got) != 1 || got[0] != want {
			t.Errorf("param %s = %v, want %q", k, got, want)
		}
	}

	if len(conns) != 1 || conns[0].From != "Bern" || conns[0].Legs[0].Line != "IC 6" {
		t.Fatalf("connections decoded wrong: %+v", conns)
	}
}

func TestFetchConnectionsDepartureOmitsTimeType(t *testing.T) {
	var gotQuery map[string][]string
	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{"connections": [{"from": "Bern"}]}`))
	})

	if _, err := FetchConnections("Bern", "Luzern", nil, "", "", false, 3); err != nil {
		t.Fatal(err)
	}

	for _, param := range []string{"time_type", "date", "time"} {
		if _, ok := gotQuery[param]; ok {
			t.Errorf("param %s should be absent, got %v", param, gotQuery[param])
		}
	}
}

func TestFetchConnectionsAPIMessage(t *testing.T) {
	withTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		// The API reports user errors as messages with HTTP 200.
		_, _ = w.Write([]byte(`{"messages": ["Stop xyz not found."], "request": null, "eof": 1}`))
	})

	_, err := FetchConnections("xyz", "Bern", nil, "", "", false, 3)
	if err == nil {
		t.Fatal("expected error from API message")
	}
	if !strings.Contains(err.Error(), "Stop xyz not found.") {
		t.Fatalf("error = %v, want it to contain the API message", err)
	}
}

func TestFetchConnectionsHTTPError(t *testing.T) {
	withTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	})

	if _, err := FetchConnections("Bern", "Luzern", nil, "", "", false, 4); err == nil {
		t.Fatal("expected error on HTTP 429")
	}
}

func TestFetchConnectionsEmptyResult(t *testing.T) {
	withTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"connections": []}`))
	})

	conns, err := FetchConnections("Bern", "Luzern", nil, "", "", false, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(conns) != 0 {
		t.Fatalf("connections = %v, want empty", conns)
	}
}
