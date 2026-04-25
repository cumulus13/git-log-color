package main

import "regexp"

type VersionMessage struct {
	Version     string
	Date        string
	Time        string
	Microsecond string
}

var versionPattern = regexp.MustCompile(`^version:\s*([\d.]+)\s*~\s*(\d{4}-\d{2}-\d{2})\s+(\d{2}:\d{2}:\d{2}):(\d+)$`)

func ParseVersionMessage(message string) *VersionMessage {
	match := versionPattern.FindStringSubmatch(message)
	if match == nil {
		return nil
	}

	return &VersionMessage{
		Version:     match[1],
		Date:        match[2],
		Time:        match[3],
		Microsecond: match[4],
	}
}
