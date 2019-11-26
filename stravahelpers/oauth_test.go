package stravahelpers

import (
	"testing"
)

// TestLoadTokens tests that the loadTokens function fails if the input secrets has values that are defaults.
func TestLoadTokens(t *testing.T) {
	var sec secrets

	_, _, err := loadTokens(sec)
	if err == nil {
		t.Error(`loadTokens did not return error on newly initialized input`)
	}
}

// TestAPICallWithBlankTokens makes a call for a Strava Activity but does not initialize any of the
// auth tokens.
func TestAPICallWithBlankTokens(t *testing.T) {
	params := map[string]uint64{
		"id": 2877607175,
	}
	_, err := StravaAPIGetResponse("StravaGetActivityPath", params)
	if err == nil {
		t.Error("the strava call should have failed")
	}
}
