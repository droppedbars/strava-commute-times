package main

import "testing"

// TestLoadTokens tests that the loadTokens function fails if the input secrets has values that are defaults.
func TestLoadTokens(t *testing.T) {
	var sec secrets

	_, _, err := loadTokens(sec)
	if err == nil {
		t.Error(`loadTokens did not return error on newly initialized input`)
	}
}
