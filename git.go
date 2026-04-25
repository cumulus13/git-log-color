package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const fieldSep = "\x1f"
const recordSep = "\x1e"

type GitClient struct {
	Bin      string
	RepoPath string
}

type GitCommit struct {
	Hash          string   `json:"hash" xml:"hash"`
	ShortHash     string   `json:"short_hash" xml:"short_hash"`
	Parents       []string `json:"parents,omitempty" xml:"parents>parent,omitempty"`
	Refs          []string `json:"refs,omitempty" xml:"refs>ref,omitempty"`
	AuthorName    string   `json:"author_name" xml:"author_name"`
	AuthorEmail   string   `json:"author_email" xml:"author_email"`
	AuthorDate    string   `json:"author_date" xml:"author_date"`
	Subject       string   `json:"subject" xml:"subject"`
	Body          string   `json:"body,omitempty" xml:"body,omitempty"`
	IsMergeCommit bool     `json:"is_merge_commit" xml:"is_merge_commit"`
	Version       string   `json:"version,omitempty" xml:"version,omitempty"`
	VersionDate   string   `json:"version_date,omitempty" xml:"version_date,omitempty"`
	VersionTime   string   `json:"version_time,omitempty" xml:"version_time,omitempty"`
	Microsecond   string   `json:"microsecond,omitempty" xml:"microsecond,omitempty"`
}

func (g GitClient) ValidateRepo() error {
	cmd := exec.Command(g.Bin, "rev-parse", "--is-inside-work-tree")
	cmd.Dir = g.RepoPath
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("not a git repository: %s", g.RepoPath)
	}
	if strings.TrimSpace(string(out)) != "true" {
		return fmt.Errorf("not a git repository: %s", g.RepoPath)
	}
	return nil
}

func (g GitClient) RunGraphLog(limit int, all bool, oneline bool, extraArgs []string) ([]string, error) {
	args := []string{"--no-pager", "log", fmt.Sprintf("-n%d", limit), "--graph", "--decorate=full"}
	if oneline {
		args = append(args, "--oneline")
	} else {
		args = append(args, "--date=iso")
	}
	if all {
		args = append(args, "--all")
	}
	args = append(args, normalizeExtraArgs(extraArgs)...)

	out, err := g.run(args...)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines, nil
}

func (g GitClient) RunStructuredLog(limit int, all bool, extraArgs []string) ([]GitCommit, error) {
	format := strings.Join([]string{
		"%H",
		"%h",
		"%P",
		"%D",
		"%an",
		"%ae",
		"%aI",
		"%s",
		"%b",
	}, fieldSep)

	args := []string{
		"--no-pager",
		"log",
		fmt.Sprintf("-n%d", limit),
		fmt.Sprintf("--pretty=format:%s%s", format, recordSep),
		"--date=iso-strict",
	}
	if all {
		args = append(args, "--all")
	}
	args = append(args, normalizeExtraArgs(extraArgs)...)

	out, err := g.run(args...)
	if err != nil {
		return nil, err
	}

	return parseStructuredLog(out)
}

func (g GitClient) run(args ...string) (string, error) {
	cmd := exec.Command(g.Bin, args...)
	cmd.Dir = g.RepoPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", errors.New(msg)
	}

	return string(out), nil
}

func normalizeExtraArgs(args []string) []string {
	if len(args) > 0 && args[0] == "--" {
		return args[1:]
	}
	return args
}

func parseStructuredLog(raw string) ([]GitCommit, error) {
	records := strings.Split(raw, recordSep)
	commits := make([]GitCommit, 0, len(records))

	for _, record := range records {
		record = strings.TrimSpace(record)
		if record == "" {
			continue
		}

		fields := strings.Split(record, fieldSep)
		if len(fields) < 9 {
			return nil, fmt.Errorf("unexpected git log record: %q", record)
		}

		commit := GitCommit{
			Hash:        fields[0],
			ShortHash:   fields[1],
			Parents:     splitNonEmpty(fields[2], " "),
			Refs:        normalizeRefs(fields[3]),
			AuthorName:  fields[4],
			AuthorEmail: fields[5],
			AuthorDate:  normalizeDate(fields[6]),
			Subject:     fields[7],
			Body:        strings.TrimSpace(fields[8]),
		}
		commit.IsMergeCommit = len(commit.Parents) > 1

		if version := ParseVersionMessage(commit.Subject); version != nil {
			commit.Version = version.Version
			commit.VersionDate = version.Date
			commit.VersionTime = version.Time
			commit.Microsecond = version.Microsecond
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

func splitNonEmpty(value string, sep string) []string {
	parts := strings.Split(strings.TrimSpace(value), sep)
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func normalizeRefs(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func normalizeDate(value string) string {
	if value == "" {
		return value
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return value
	}
	return t.Format(time.RFC3339)
}
