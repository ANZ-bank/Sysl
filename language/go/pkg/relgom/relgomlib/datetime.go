package relgomlib

import (
	"regexp"
	"time"
)

var ymdRE = regexp.MustCompile(`\A\d{4}-\d{2}-\d{2}\z`)

type DateTimeString string

func NewDateTimeString(t *time.Time) *DateTimeString {
	if t == nil {
		return nil
	}
	result := DateTimeString(t.Format(time.RFC3339))
	return &result
}

func (s *DateTimeString) Unstage() (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	if ymdRE.MatchString(string(*s)) {
		t, err := time.Parse(time.RFC3339, string(*s)+"T00:00:00Z")
		return &t, err
	}
	t, err := time.Parse(time.RFC3339, string(*s))
	return &t, err
}
