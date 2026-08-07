package main

import "strings"

// defaultAuthorPalette is the ordered list of styles assigned to distinct
// commit authors. Order matters: the first author encountered gets the
// first color, the second author gets the second color, and so on. This
// keeps colors deterministic and sequential (not random), while still
// giving each unique author a consistently distinct color across a run.
// If there are more authors than palette entries, the palette wraps
// around and repeats from the start.
//
// Each entry follows the same shape as the original default "author"
// style (white italic text on a colored background) - only the
// background color changes between owners, so the "author block" look
// is preserved.
func defaultAuthorPalette() []string {
	return []string{
		"white bg:#AA0000 italic", // red
		"white bg:#006400 italic", // green
		"white bg:#00008B italic", // blue
		"white bg:#8B5A00 italic", // brown/gold
		"white bg:#8B008B italic", // magenta
		"white bg:#008B8B italic", // cyan
		"white bg:#CC5500 italic", // orange
		"white bg:#4B0082 italic", // indigo
		"white bg:#006666 italic", // teal
		"white bg:#99004C italic", // pink/maroon
		"white bg:#556B2F italic", // olive
		"white bg:#2F4F4F italic", // slate
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
