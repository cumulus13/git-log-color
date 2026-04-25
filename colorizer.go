package main

import (
	"fmt"
	"regexp"
	"strings"
)

type TextRenderer struct {
	styler Styler
	icons  IconSet
}

func NewTextRenderer(cfg Config) TextRenderer {
	return TextRenderer{
		styler: NewStyler(cfg),
		icons:  NewIconSet(cfg),
	}
}

var graphLinePattern = regexp.MustCompile(`^([*|/\\\s\-_+]+)(.*)$`)
var commitLinePattern = regexp.MustCompile(`^commit ([0-9a-f]+)(.*)$`)
var authorLinePattern = regexp.MustCompile(`^Author:\s*([^<]+)<([^>]+)>`)
var dateLinePattern = regexp.MustCompile(`^Date:\s*(.+)$`)
var onelinePattern = regexp.MustCompile(`^([*|/\\\s\-_+]+)([0-9a-f]+)\s*(.*)$`)

func (r TextRenderer) RenderLine(line string, oneline bool) string {
	if strings.TrimSpace(line) == "" {
		return ""
	}

	if oneline {
		return r.renderOneLine(line)
	}
	return r.renderFullLine(line)
}

func (r TextRenderer) renderFullLine(line string) string {
	match := graphLinePattern.FindStringSubmatch(line)
	if match == nil {
		return r.styler.Apply("message", line)
	}

	graphPart := match[1]
	content := strings.TrimSpace(match[2])
	graph := r.colorizeGraph(graphPart)

	if content == "" {
		return graph
	}

	if commitMatch := commitLinePattern.FindStringSubmatch(content); commitMatch != nil {
		hash := commitMatch[1]
		refs := strings.TrimSpace(commitMatch[2])

		var out strings.Builder
		out.WriteString(graph)
		out.WriteString(r.styler.Apply("commit_graph", r.icons.Get("commit")+"commit "))
		out.WriteString(r.styler.Apply("commit", shortenHash(hash)))
		if refs != "" {
			out.WriteString(r.colorizeRefs(refs))
		}
		return out.String()
	}

	if authorMatch := authorLinePattern.FindStringSubmatch(content); authorMatch != nil {
		return graph +
			r.styler.Apply("username", r.icons.Get("username")+" Author: ") +
			r.styler.Apply("author", strings.TrimSpace(authorMatch[1])) +
			" " +
			r.styler.Apply("remote", r.icons.Get("email")+"<"+authorMatch[2]+">")
	}

	if dateMatch := dateLinePattern.FindStringSubmatch(content); dateMatch != nil {
		return graph +
			r.styler.Apply("date_string", r.icons.Get("date")+" Date: ") +
			r.styler.Apply("date", dateMatch[1])
	}

	if version := ParseVersionMessage(content); version != nil {
		return graph +
			r.styler.Apply("message", r.icons.Get("version")+" version:") +
			r.styler.Apply("version_number", " "+version.Version) +
			r.styler.Apply("date_number", " ~ "+version.Date) +
			r.styler.Apply("time_number", " "+version.Time) +
			r.styler.Apply("microsecond", ":"+version.Microsecond)
	}

	return graph + r.styler.Apply("message", r.icons.Get("version")+content)
}

func (r TextRenderer) renderOneLine(line string) string {
	match := onelinePattern.FindStringSubmatch(line)
	if match == nil {
		return r.styler.Apply("message", line)
	}

	graph := r.colorizeGraph(match[1])
	hash := r.styler.Apply("commit", shortenHash(match[2]))
	messageAndRefs := strings.TrimSpace(match[3])

	if messageAndRefs == "" {
		return graph + hash
	}

	if strings.Contains(messageAndRefs, "(") && strings.Contains(messageAndRefs, ")") {
		return graph + hash + " " + r.colorizeRefs(messageAndRefs)
	}

	return graph + hash + " " + r.styler.Apply("message", messageAndRefs)
}

func (r TextRenderer) colorizeGraph(graph string) string {
	var out strings.Builder

	for _, char := range graph {
		switch char {
		case '*':
			out.WriteString(r.styler.Apply("commit_dot", "*"))
		case '|':
			out.WriteString(r.styler.Apply("graph_main", "│"))
		case '/', '\\':
			out.WriteString(r.styler.Apply("graph_merge", string(char)))
		case '-', '_':
			out.WriteString(r.styler.Apply("graph_branch", string(char)))
		case '+':
			out.WriteString(r.styler.Apply("merge_dot", "+"))
		default:
			out.WriteString(r.styler.Apply("graph_dim", string(char)))
		}
	}

	return out.String()
}

func (r TextRenderer) colorizeRefs(input string) string {
	match := regexp.MustCompile(`\(([^)]+)\)`).FindStringSubmatchIndex(input)
	if match == nil {
		return r.styler.Apply("message", " "+input)
	}

	before := input[:match[0]]
	content := input[match[2]:match[3]]
	after := input[match[1]:]

	var out strings.Builder
	if strings.TrimSpace(before) != "" {
		out.WriteString(" ")
		out.WriteString(r.styler.Apply("message", before))
	}

	out.WriteString(" (")
	parts := strings.Split(content, ",")
	for index, part := range parts {
		part = strings.TrimSpace(part)

		switch {
		case strings.Contains(part, "HEAD"), strings.Contains(part, "master"), strings.Contains(part, "main"):
			out.WriteString(r.styler.Apply("branch", part))
		case strings.Contains(part, "tag:"):
			out.WriteString(r.styler.Apply("tag", fmt.Sprintf("%s%s", r.icons.Get("tag"), part)))
		case strings.Contains(part, "origin"), strings.Contains(part, "upstream"):
			out.WriteString(r.styler.Apply("remote", part))
		default:
			out.WriteString(r.styler.Apply("graph_dim", part))
		}

		if index < len(parts)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")

	if strings.TrimSpace(after) != "" {
		out.WriteString(r.styler.Apply("message", after))
	}

	return out.String()
}

func shortenHash(value string) string {
	if len(value) <= 7 {
		return value
	}
	return value[:7]
}
