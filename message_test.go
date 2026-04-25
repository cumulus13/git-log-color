package main

import "testing"

func TestParseVersionMessage(t *testing.T) {
	version := ParseVersionMessage("version: 1.11.34 ~ 2025-09-23 00:55:42:978567")
	if version == nil {
		t.Fatal("expected version message to parse")
	}

	if version.Version != "1.11.34" {
		t.Fatalf("unexpected version: %s", version.Version)
	}
	if version.Date != "2025-09-23" {
		t.Fatalf("unexpected date: %s", version.Date)
	}
	if version.Time != "00:55:42" {
		t.Fatalf("unexpected time: %s", version.Time)
	}
	if version.Microsecond != "978567" {
		t.Fatalf("unexpected microsecond: %s", version.Microsecond)
	}
}

func TestParseVersionMessageRejectsOtherMessages(t *testing.T) {
	if ParseVersionMessage("feat: add formatter") != nil {
		t.Fatal("expected regular subject not to parse as version message")
	}
}
