package controllers

import "testing"

func TestLinkedinURL(t *testing.T) {
	testCases := []struct {
		url   string
		match bool
	}{
		{"https://www.linkedin.com/in/james-bond-007/", true},
		{"https://linkedin.com/in/marllos-p-a383641b2/", true},
		{"https://", false},
		{"https://www.linkedin.com/in/", false},
	}

	for _, tc := range testCases {
		got, err := isValidLinkedIn(tc.url)
		if err != nil {
			t.Errorf("isValidLinkedIn(%s) returned unexpected error: %v", tc.url, err)
		}
		if got != tc.match {
			t.Errorf("isValidLinkedIn(%s) = %t, want %t", tc.url, got, tc.match)
		}
	}
}
