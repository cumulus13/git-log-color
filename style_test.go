package main

import (
	"strings"
	"testing"
)

func TestParseStyleSpec(t *testing.T) {
	ansi, err := parseStyleSpec("bold #FFAAFF bg:#001122 italic")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, part := range []string{"1", "3", "38;2;255;170;255", "48;2;0;17;34"} {
		if !strings.Contains(ansi, part) {
			t.Fatalf("expected ANSI sequence to contain %q, got %q", part, ansi)
		}
	}
}

func TestSupportedFormats(t *testing.T) {
	for _, value := range []string{"text", "json", "table", "xml", "yaml", "toml"} {
		if !isSupportedFormat(value) {
			t.Fatalf("expected format %q to be supported", value)
		}
	}
}
