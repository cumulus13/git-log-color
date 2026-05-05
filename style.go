package main

import (
	"fmt"
	"strconv"
	"strings"
)

const ansiReset = "\033[0m"

func defaultStyles() map[string]string {
	return map[string]string{
		"commit":         "bold yellow",
		"branch":         "bold cyan",
		"tag":            "bold #FFAAFF",
		"tag_version":    "white bg:#0000FF italic",
		"remote":         "green",
		"author":         "white bg:#AA0000 italic",
		"date":           "bold #AA55FF",
		"message":        "white",
		"graph_main":     "bold cyan",
		"graph_merge":    "bold green",
		"graph_branch":   "bold yellow",
		"graph_dim":      "dim white",
		"commit_dot":     "bold red",
		"merge_dot":      "bold green",
		"commit_graph":   "bold #55FF00",
		"username":       "bold #FF5500",
		"date_string":    "bold #0000FF",
		"date_number":    "bold #AAFFFF",
		"version_string": "bold #2CFF2C",
		"version_number": "bold #FFFF00",
		"time_number":    "bold #55AAFF",
		"microsecond":    "bold #AA0000",
	}
}

type Styler struct {
	styles  map[string]string
	enabled bool
}

func NewStyler(cfg Config) Styler {
	return Styler{
		styles:  cfg.Styles,
		enabled: cfg.UI.ColorEnabled,
	}
}

func (s Styler) Apply(name string, text string) string {
	if !s.enabled {
		return text
	}

	spec, ok := s.styles[name]
	if !ok || strings.TrimSpace(spec) == "" {
		return text
	}

	ansi, err := parseStyleSpec(spec)
	if err != nil {
		return text
	}

	return ansi + text + ansiReset
}

func parseStyleSpec(spec string) (string, error) {
	tokens := strings.Fields(strings.TrimSpace(spec))
	if len(tokens) == 0 {
		return "", nil
	}

	var codes []string
	for _, token := range tokens {
		switch {
		case token == "bold":
			codes = append(codes, "1")
		case token == "dim":
			codes = append(codes, "2")
		case token == "italic":
			codes = append(codes, "3")
		case strings.HasPrefix(token, "bg:"):
			rgb, err := parseColor(token[3:])
			if err != nil {
				return "", err
			}
			codes = append(codes, fmt.Sprintf("48;2;%d;%d;%d", rgb[0], rgb[1], rgb[2]))
		default:
			rgb, err := parseColor(token)
			if err != nil {
				return "", err
			}
			codes = append(codes, fmt.Sprintf("38;2;%d;%d;%d", rgb[0], rgb[1], rgb[2]))
		}
	}

	return "\033[" + strings.Join(codes, ";") + "m", nil
}

func parseColor(value string) ([3]int, error) {
	named := map[string][3]int{
		"black":   {0, 0, 0},
		"red":     {255, 0, 0},
		"green":   {0, 170, 0},
		"yellow":  {255, 215, 0},
		"blue":    {0, 102, 255},
		"magenta": {255, 0, 255},
		"cyan":    {0, 255, 255},
		"white":   {255, 255, 255},
		"gray":    {160, 160, 160},
		"grey":    {160, 160, 160},
	}

	if rgb, ok := named[strings.ToLower(value)]; ok {
		return rgb, nil
	}

	if strings.HasPrefix(value, "#") && len(value) == 7 {
		r, err := strconv.ParseInt(value[1:3], 16, 64)
		if err != nil {
			return [3]int{}, err
		}
		g, err := strconv.ParseInt(value[3:5], 16, 64)
		if err != nil {
			return [3]int{}, err
		}
		b, err := strconv.ParseInt(value[5:7], 16, 64)
		if err != nil {
			return [3]int{}, err
		}
		return [3]int{int(r), int(g), int(b)}, nil
	}

	return [3]int{}, fmt.Errorf("unsupported color %q", value)
}
