package main

import (
	"context"
	"os"
	"testing"
)

func TestFilterURL(t *testing.T) {
	t.Setenv("LINK_REGEXP", `^https://([\w.-]+)`) // capture domain

	tests := []struct {
		input      string
		expectOk   bool
		expectURL  string
		expectText string
	}{
		{"https://example.com/path", true, "https://example.com/path", "example.com"},
		{"http://example.com/path", false, "", ""},
		{"https://sub.domain.com/", true, "https://sub.domain.com/", "sub.domain.com"},
		{"not-a-url", false, "", ""},
	}
	for _, tt := range tests {
		url, text, ok := filterURL(tt.input)
		if ok != tt.expectOk || url != tt.expectURL || text != tt.expectText {
			t.Errorf(
				"filterURL(%q) = (%q, %q, %v), want (%q, %q, %v)",
				tt.input,
				url,
				text,
				ok,
				tt.expectURL,
				tt.expectText,
				tt.expectOk,
			)
		}
	}
}

func TestDataEnvVars(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	os.Unsetenv("UPSUN_API_TOKEN")
	os.Unsetenv("UPSUN_PROJECT_ID")
	// This should fatally log, but we can't catch log.Fatalln easily in tests.
	// So just check that it panics or exits.
	_, err := data(ctx)
	if err == nil {
		t.Error("expected error when env vars are missing, got nil")
	}
}
