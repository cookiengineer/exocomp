package colors

import (
	utils_tty "exocomp/utils/tty"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestThemeCSSMatchesPalette(t *testing.T) {

	data, err0 := os.ReadFile("../../../../app/public/app/theme.css")

	if err0 != nil {
		t.Fatalf("Read theme.css: %s", err0.Error())
	}

	values := make(map[string]string)

	pattern := regexp.MustCompile(`(?m)^\s*--([a-z0-9-]+):\s*(#[0-9a-fA-F]{6})\s*;`)

	for _, match := range pattern.FindAllStringSubmatch(string(data), -1) {
		values[match[1]] = strings.ToLower(match[2])
	}

	checks := map[string]utils_tty.Color{
		"main-color":             MainColor,
		"main-background":        MainBackground,
		"text-muted":             TextMuted,
		"background-raised":      BackgroundRaised,
		"background-raised-high": BackgroundRaisedHigh,
		"background-raised-max":  BackgroundRaisedMax,
		"background-composer":    BackgroundComposer,
		"border-base":            BorderBase,
		"scrollbar-thumb":        ScrollbarThumb,
		"accent":                 Accent,
		"accent-text":            AccentText,
		"feedback-error":         FeedbackError,
		"feedback-warning":       FeedbackWarning,
		"feedback-success":       FeedbackSuccess,
		"feedback-info":          FeedbackInfo,
		"feedback-thinking":      FeedbackThinking,
		"markdown-block-quote":   MarkdownBlockQuote,
		"markdown-emphasis":      MarkdownEmphasis,
		"diff-text-added":        DiffTextAdded,
		"diff-text-removed":      DiffTextRemoved,
		"diff-text-context":      DiffTextContext,
		"diff-text-hunk-header":  DiffTextHunkHeader,
		"diff-background-added":  DiffBackgroundAdded,
		"diff-highlight-added":   DiffHighlightAdded,
		"diff-highlight-removed": DiffHighlightRemoved,
		"syntax-type":            SyntaxType,
	}

	for name, color := range checks {

		value, ok := values[name]

		if ok == false {
			t.Errorf("theme.css is missing --%s", name)
			continue
		}

		expected := strings.ToLower(color.Hex())

		if value != expected {
			t.Errorf("--%s is %s, want %s", name, value, expected)
		}

	}

}
