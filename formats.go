package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type xmlCommitList struct {
	XMLName xml.Name    `xml:"commits"`
	Commits []GitCommit `xml:"commit"`
}

func isSupportedFormat(value string) bool {
	switch value {
	case "text", "json", "table", "xml", "yaml", "toml":
		return true
	default:
		return false
	}
}

func FormatCommits(format string, commits []GitCommit, cfg Config) (string, error) {
	switch format {
	case "json":
		data, err := json.MarshalIndent(commits, "", "  ")
		return string(data), err
	case "xml":
		data, err := xml.MarshalIndent(xmlCommitList{Commits: commits}, "", "  ")
		if err != nil {
			return "", err
		}
		return xml.Header + string(data), nil
	case "yaml":
		return toYAML(commits), nil
	case "toml":
		return toTOML(commits), nil
	case "table":
		return toTable(commits, cfg), nil
	default:
		return "", fmt.Errorf("unsupported output format %q", format)
	}
}

func toTable(commits []GitCommit, cfg Config) string {
	icons := NewIconSet(cfg)
	rows := [][]string{
		{"HASH", "AUTHOR", "DATE", "MESSAGE", "REFS"},
	}

	for _, commit := range commits {
		message := commit.Subject
		if commit.Version != "" {
			message = icons.Get("version") + "version: " + commit.Version + " ~ " + commit.VersionDate + " " + commit.VersionTime + ":" + commit.Microsecond
		}
		if len(commit.Refs) == 0 {
			rows = append(rows, []string{
				commit.ShortHash,
				commit.AuthorName,
				commit.AuthorDate,
				message,
				"",
			})
			continue
		}

		rows = append(rows, []string{
			commit.ShortHash,
			commit.AuthorName,
			commit.AuthorDate,
			message,
			strings.Join(commit.Refs, ", "),
		})
	}

	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var out strings.Builder
	for rowIndex, row := range rows {
		if rowIndex == 1 {
			out.WriteString(renderTableDivider(widths))
		}
		out.WriteString(renderTableRow(row, widths))
	}
	return out.String()
}

func renderTableDivider(widths []int) string {
	var out strings.Builder
	for _, width := range widths {
		out.WriteString("+")
		out.WriteString(strings.Repeat("-", width+2))
	}
	out.WriteString("+\n")
	return out.String()
}

func renderTableRow(row []string, widths []int) string {
	var out strings.Builder
	for i, cell := range row {
		out.WriteString("| ")
		out.WriteString(cell)
		out.WriteString(strings.Repeat(" ", widths[i]-len(cell)+1))
	}
	out.WriteString("|\n")
	return out.String()
}

func toYAML(commits []GitCommit) string {
	var out strings.Builder
	for _, commit := range commits {
		out.WriteString("- hash: " + yamlScalar(commit.Hash) + "\n")
		out.WriteString("  short_hash: " + yamlScalar(commit.ShortHash) + "\n")
		out.WriteString("  author_name: " + yamlScalar(commit.AuthorName) + "\n")
		out.WriteString("  author_email: " + yamlScalar(commit.AuthorEmail) + "\n")
		out.WriteString("  author_date: " + yamlScalar(commit.AuthorDate) + "\n")
		out.WriteString("  subject: " + yamlScalar(commit.Subject) + "\n")
		if commit.Body != "" {
			out.WriteString("  body: " + yamlScalar(commit.Body) + "\n")
		}
		out.WriteString("  is_merge_commit: " + strconv.FormatBool(commit.IsMergeCommit) + "\n")
		out.WriteString("  parents:\n")
		for _, parent := range commit.Parents {
			out.WriteString("    - " + yamlScalar(parent) + "\n")
		}
		out.WriteString("  refs:\n")
		for _, ref := range commit.Refs {
			out.WriteString("    - " + yamlScalar(ref) + "\n")
		}
		if commit.Version != "" {
			out.WriteString("  version: " + yamlScalar(commit.Version) + "\n")
			out.WriteString("  version_date: " + yamlScalar(commit.VersionDate) + "\n")
			out.WriteString("  version_time: " + yamlScalar(commit.VersionTime) + "\n")
			out.WriteString("  microsecond: " + yamlScalar(commit.Microsecond) + "\n")
		}
	}
	return strings.TrimRight(out.String(), "\n")
}

func yamlScalar(value string) string {
	escaped := strings.ReplaceAll(value, `"`, `\"`)
	escaped = strings.ReplaceAll(escaped, "\n", `\n`)
	return `"` + escaped + `"`
}

func toTOML(commits []GitCommit) string {
	var out strings.Builder
	for _, commit := range commits {
		out.WriteString("[[commit]]\n")
		out.WriteString("hash = " + tomlString(commit.Hash) + "\n")
		out.WriteString("short_hash = " + tomlString(commit.ShortHash) + "\n")
		out.WriteString("author_name = " + tomlString(commit.AuthorName) + "\n")
		out.WriteString("author_email = " + tomlString(commit.AuthorEmail) + "\n")
		out.WriteString("author_date = " + tomlString(commit.AuthorDate) + "\n")
		out.WriteString("subject = " + tomlString(commit.Subject) + "\n")
		out.WriteString("body = " + tomlString(commit.Body) + "\n")
		out.WriteString("is_merge_commit = " + strconv.FormatBool(commit.IsMergeCommit) + "\n")
		out.WriteString("parents = " + tomlArray(commit.Parents) + "\n")
		out.WriteString("refs = " + tomlArray(commit.Refs) + "\n")
		if commit.Version != "" {
			out.WriteString("version = " + tomlString(commit.Version) + "\n")
			out.WriteString("version_date = " + tomlString(commit.VersionDate) + "\n")
			out.WriteString("version_time = " + tomlString(commit.VersionTime) + "\n")
			out.WriteString("microsecond = " + tomlString(commit.Microsecond) + "\n")
		}
		out.WriteString("\n")
	}
	return strings.TrimSpace(out.String())
}

func tomlString(value string) string {
	return strconv.Quote(value)
}

func tomlArray(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, tomlString(value))
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func MarshalConfigExample(cfg Config) string {
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(cfg)
	return strings.TrimSpace(out.String())
}
