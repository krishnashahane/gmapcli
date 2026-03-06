package terminal

import (
	"os"
	"strings"
)

// Palette controls ANSI color output.
type Palette struct {
	active bool
}

// NewPalette creates a color helper.
func NewPalette(enabled bool) Palette {
	return Palette{active: enabled}
}

func (p Palette) Strong(s string) string  { return p.esc("1", s) }
func (p Palette) Teal(s string) string    { return p.esc("36", s) }
func (p Palette) Lime(s string) string    { return p.esc("32", s) }
func (p Palette) Amber(s string) string   { return p.esc("33", s) }
func (p Palette) Muted(s string) string   { return p.esc("2", s) }

func (p Palette) esc(code, s string) string {
	if !p.active {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

// ShouldColor decides if ANSI codes should be emitted.
func ShouldColor(forceOff bool) bool {
	if forceOff {
		return false
	}
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	t := strings.TrimSpace(os.Getenv("TERM"))
	return t != "" && t != "dumb"
}
