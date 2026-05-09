// Package model defines the data types decoded from the transport API.
package model

import (
	"strings"
	"time"
)

// SwissLocation is the Europe/Zurich time zone, with a fixed-offset
// fallback when the tzdata lookup fails. The fallback does not handle
// the CET-CEST transition; main.go embeds time/tzdata to keep that
// path off the hot road.
var SwissLocation = func() *time.Location {
	loc, err := time.LoadLocation("Europe/Zurich")
	if err != nil {
		loc = time.FixedZone("CET", 1*60*60)
	}
	return loc
}()

// Timestamp is a time.Time decoded from the API's RFC3339-ish format.
type Timestamp struct {
	time.Time
}

// Sub returns the duration between two Timestamps.
func (t Timestamp) Sub(other Timestamp) time.Duration {
	return t.Time.Sub(other.Time)
}

// Local returns the timestamp in Swiss time, overriding time.Time.Local.
func (t Timestamp) Local() time.Time {
	return t.In(SwissLocation)
}

// UnmarshalJSON parses the API's "2006-01-02T15:04:05-0700" format.
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

// Coordinate is a geographic point.
type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Station is a transport stop with a name and coordinates.
type Station struct {
	Name       string     `json:"name"`
	Coordinate Coordinate `json:"coordinate"`
}

// Departure describes when and where a section leaves.
type Departure struct {
	Station   Station   `json:"station"`
	Scheduled Timestamp `json:"departure"`
	Platform  string    `json:"platform"`
	Delay     int       `json:"delay"`
}

// Arrival describes when and where a section reaches its end.
type Arrival struct {
	Station   Station   `json:"station"`
	Scheduled Timestamp `json:"arrival"`
	Platform  string    `json:"platform"`
	Delay     int       `json:"delay"`
}

// Section is a single leg of a connection: a vehicle journey or a walk.
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

// Connection is the full route returned by the API for one query result.
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
