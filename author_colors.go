package main

import "strings"

// defaultAuthorPalette is the ordered list of styles assigned to distinct
// commit authors. Order matters: the first author encountered gets the
// first color, the second author gets the second color, and so on. This
// keeps colors deterministic and sequential (not random), while still
// giving each unique author a consistently distinct color across a run.
// If there are more authors than palette entries, the palette wraps
// around and repeats from the start.
func defaultAuthorPalette() []string {
	return []string{
		"bold #FF5555", // red
		"bold #55FF55", // green
		"bold #5599FF", // blue
		"bold #FFD700", // gold
		"bold #FF55FF", // magenta
		"bold #55FFFF", // cyan
		"bold #FFA500", // orange
		"bold #AA55FF", // purple
		"bold #00CED1", // dark turquoise
		"bold #FF69B4", // hot pink
		"bold #9ACD32", // yellow-green
		"bold #40E0D0", // turquoise
	}
}

// AuthorColorAssigner hands out a style from a fixed palette to each
// distinct author, in the order authors are first seen. The same author
// always gets the same style for the lifetime of the assigner; different
// authors get different styles, cycling through the palette sequentially
// once every slot has been used.
type AuthorColorAssigner struct {
	palette []string
	assign  map[string]string
	order   []string
}

func NewAuthorColorAssigner(palette []string) *AuthorColorAssigner {
	if len(palette) == 0 {
		palette = defaultAuthorPalette()
	}
	return &AuthorColorAssigner{
		palette: palette,
		assign:  make(map[string]string),
	}
}

// authorKey normalizes an author name so trivial formatting differences
// (extra whitespace, case) don't accidentally split one person into two
// colors.
func authorKey(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// StyleFor returns the style spec for the given author, assigning the
// next sequential palette slot the first time this author is seen.
func (a *AuthorColorAssigner) StyleFor(name string) string {
	key := authorKey(name)
	if key == "" {
		return ""
	}
	if style, ok := a.assign[key]; ok {
		return style
	}

	index := len(a.order) % len(a.palette)
	style := a.palette[index]
	a.assign[key] = style
	a.order = append(a.order, key)
	return style
}

// Apply styles the given text as the given author's commit, using that
// author's sequentially-assigned color.
func (a *AuthorColorAssigner) Apply(styler Styler, name string, text string) string {
	style := a.StyleFor(name)
	if style == "" {
		return styler.Apply("author", text)
	}
	return styler.ApplySpec(style, text)
}
