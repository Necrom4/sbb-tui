// Package model defines data types for the Swiss public transport API responses.
package model

import (
	"strings"
	"time"
)

// SwissLocation is the time zone for Switzerland (CET/CEST, Europe/Zurich).
// main.go imports _ "time/tzdata" so the embedded TZ database is always
// available; the fixed-offset fallback is a last-resort safety net only and
// does not handle the CET↔CEST daylight-saving transition.
var SwissLocation = func() *time.Location {
	loc, err := time.LoadLocation("Europe/Zurich")
	if err != nil {
		loc = time.FixedZone("CET", 1*60*60)
	}
	return loc
}()

// Timestamp wraps time.Time with custom JSON unmarshaling for the SBB API
// date format (2006-01-02T15:04:05-0700).
type Timestamp struct {
	time.Time
}

// Sub returns the duration between two Timestamps.
func (t Timestamp) Sub(other Timestamp) time.Duration {
	return t.Time.Sub(other.Time)
}

// Local overrides time.Time.Local() to return the timestamp in Swiss time
// (CET/CEST, Europe/Zurich) instead of the system's local timezone.
func (t Timestamp) Local() time.Time {
	return t.Time.In(SwissLocation)
}

// UnmarshalJSON parses the SBB API date format into a Timestamp.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	parsed, err := time.Parse("2006-01-02T15:04:05-0700", s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// Coordinate represents a geographic position.
type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Station represents a transport stop or station.
type Station struct {
	Name       string     `json:"name"`
	Coordinate Coordinate `json:"coordinate"`
}

// Departure holds departure information for a connection section.
type Departure struct {
	Station   Station   `json:"station"`
	Scheduled Timestamp `json:"departure"`
	Platform  string    `json:"platform"`
	Delay     int       `json:"delay"`
}

// Arrival holds arrival information for a connection section.
type Arrival struct {
	Station   Station   `json:"station"`
	Scheduled Timestamp `json:"arrival"`
	Platform  string    `json:"platform"`
	Delay     int       `json:"delay"`
}

// Section represents a leg of a connection (either a journey or a walk).
type Section struct {
	Journey *struct {
		Category string `json:"category"`
		Number   string `json:"number"`
		Operator string `json:"operator"`
		To       string `json:"to"`
	} `json:"journey"`
	Walk *struct {
		Duration  int       `json:"duration"`
		Departure Departure `json:"departure"`
		Arrival   Arrival   `json:"arrival"`
	} `json:"walk"`
	Departure Departure `json:"departure"`
	Arrival   Arrival   `json:"arrival"`
}

// Connection represents a complete route between two stations.
type Connection struct {
	From struct {
		Station   Station   `json:"station"`
		Departure Timestamp `json:"departure"`
		Delay     int       `json:"delay"`
		Platform  string    `json:"platform"`
	} `json:"from"`

	To struct {
		Station  Station   `json:"station"`
		Arrival  Timestamp `json:"arrival"`
		Platform string    `json:"platform"`
	} `json:"to"`

	Duration  string    `json:"duration"`
	Transfers int       `json:"transfers"`
	Sections  []Section `json:"sections"`
}
