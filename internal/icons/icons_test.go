package icons_test

import (
	"testing"

	"github.com/mohokh67/portly/internal/icons"
)

func TestResolveKnownProcess(t *testing.T) {
	icon := icons.Resolve("node", icons.Emoji)
	if icon == "" {
		t.Error("expected non-empty icon for 'node'")
	}
	if icon != "🟢" {
		t.Errorf("expected 🟢 for node emoji, got %q", icon)
	}
}

func TestResolveUnknownProcess(t *testing.T) {
	icon := icons.Resolve("somethingobscure", icons.Emoji)
	if icon == "" {
		t.Error("expected fallback icon for unknown process")
	}
}

func TestResolveNoneStyle(t *testing.T) {
	icon := icons.Resolve("node", icons.None)
	if icon != "" {
		t.Errorf("expected empty string for None style, got %q", icon)
	}
}

func TestResolvePostgresVariants(t *testing.T) {
	for _, name := range []string{"postgres", "postgresql", "pg"} {
		icon := icons.Resolve(name, icons.Emoji)
		if icon != "🐘" {
			t.Errorf("expected 🐘 for %q, got %q", name, icon)
		}
	}
}

func TestDetectStyleNerdFontsEnv(t *testing.T) {
	t.Setenv("NERD_FONTS", "1")
	t.Setenv("TERM_PROGRAM", "")
	style := icons.DetectStyle()
	if style != icons.NerdFont {
		t.Errorf("expected NerdFont when NERD_FONTS=1, got %v", style)
	}
}

func TestDetectStyleKnownTerminal(t *testing.T) {
	t.Setenv("NERD_FONTS", "")
	t.Setenv("TERM_PROGRAM", "iTerm.app")
	style := icons.DetectStyle()
	if style != icons.NerdFont {
		t.Errorf("expected NerdFont for iTerm.app, got %v", style)
	}
}

func TestDetectStyleUnknownTerminalFallsBackToEmoji(t *testing.T) {
	t.Setenv("NERD_FONTS", "")
	t.Setenv("TERM_PROGRAM", "SomeRandomTerm")
	style := icons.DetectStyle()
	if style != icons.Emoji {
		t.Errorf("expected Emoji fallback for unknown terminal, got %v", style)
	}
}

func TestParseStyle(t *testing.T) {
	cases := map[string]icons.IconStyle{
		"nerdfont": icons.NerdFont,
		"emoji":    icons.Emoji,
		"none":     icons.None,
		"auto":     icons.Auto,
		"unknown":  icons.Auto,
	}
	for input, want := range cases {
		got := icons.ParseStyle(input)
		if got != want {
			t.Errorf("ParseStyle(%q) = %v, want %v", input, got, want)
		}
	}
}
