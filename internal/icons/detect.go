package icons

import "os"

// nerdFontTerminals is the list of $TERM_PROGRAM values known to support Nerd Fonts.
var nerdFontTerminals = map[string]bool{
	"iTerm.app": true,
	"WezTerm":   true,
	"ghostty":   true,
	"kitty":     true,
	"Alacritty": true,
}

// DetectStyle determines the best icon style for the current terminal.
// Order: $NERD_FONTS env var → $TERM_PROGRAM → emoji fallback.
func DetectStyle() IconStyle {
	if os.Getenv("NERD_FONTS") == "1" {
		return NerdFont
	}
	if nerdFontTerminals[os.Getenv("TERM_PROGRAM")] {
		return NerdFont
	}
	return Emoji
}

// ParseStyle converts a --icons flag string to IconStyle.
func ParseStyle(s string) IconStyle {
	switch s {
	case "nerdfont":
		return NerdFont
	case "emoji":
		return Emoji
	case "none":
		return None
	default:
		return Auto
	}
}

// ResolveAuto returns the icon with auto-detected style.
func ResolveAuto(processName string) string {
	return Resolve(processName, DetectStyle())
}
