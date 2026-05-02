// Package style centralises the Optiqor CLI's visual language: colors,
// badges, dividers, and section formatters. Renderers compose these
// styles; they never reach for raw ANSI codes.
//
// All styles auto-degrade based on terminal capability: when output is
// not a TTY, when NO_COLOR is set, or when the user passes --no-color,
// every Style here renders as its plain-text equivalent.
package style

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// BrandGlyph is the ASCII stand-in for the optiqor logomark — a
// circular Q rendered at terminal scale. Using a single glyph keeps the
// header readable on every emulator (we don't depend on Nerd Fonts or
// Unicode powerline glyphs, which still render as boxes on stock CI).
const BrandGlyph = "◐"

// Theme bundles the entire palette + reusable styles. Construct one
// per render invocation via NewTheme so colour-vs-plain is a single
// decision, not threaded through every helper.
type Theme struct {
	UseColor bool

	// Brand
	Brand        lipgloss.Style
	BrandMark    lipgloss.Style
	Tagline      lipgloss.Style
	HeaderBorder lipgloss.Style

	// Boxed-finding card + signal-bar palette
	CardBorder  lipgloss.Style
	BarFilled   lipgloss.Style
	BarEmpty    lipgloss.Style
	BarOverflow lipgloss.Style // for ratios > 1 (limit < request, etc.)

	// Sections
	SectionPrimary lipgloss.Style // headline section (Cost optimizations)
	SectionBonus   lipgloss.Style // bonus section (Security)
	SectionSubtle  lipgloss.Style // light explanatory line under a section

	// Severity badges
	SevHigh lipgloss.Style
	SevMed  lipgloss.Style
	SevLow  lipgloss.Style
	SevInfo lipgloss.Style

	// Confidence
	ConfHigh lipgloss.Style
	ConfMed  lipgloss.Style
	ConfLow  lipgloss.Style

	// Output elements
	Workload   lipgloss.Style
	Title      lipgloss.Style
	Detail     lipgloss.Style
	Savings    lipgloss.Style
	BigSavings lipgloss.Style // hero number in the executive summary
	NoSavings  lipgloss.Style
	Muted      lipgloss.Style
	Divider    lipgloss.Style
	Disclosure lipgloss.Style
	CallToLink lipgloss.Style
	OK         lipgloss.Style
}

// NewTheme builds a theme. If useColor is false, every style falls back
// to plain text and bold/foreground attributes are no-ops.
//
// When useColor is true, the theme uses its own renderer with the color
// profile pinned to TrueColor so output is consistent regardless of
// what TTY detection says — the CLI's outer layer already gated on
// TTY/NO_COLOR/--no-color before deciding to call NewTheme(true).
func NewTheme(useColor bool) Theme {
	if !useColor {
		return plainTheme()
	}

	// Use our own renderer so tests and pipe-redirected output still
	// emit ANSI when the caller explicitly asked for color. The CLI's
	// outer layer is the source of truth for "should I color?".
	r := lipgloss.NewRenderer(io.Discard)
	r.SetColorProfile(termenv.TrueColor)

	// Adaptive colors: pick a value that reads well on both dark and
	// light backgrounds. Dark variants are tuned for the most common
	// terminal default (dark).
	brand := lipgloss.AdaptiveColor{Light: "#5C2EE5", Dark: "#A78BFA"}
	red := lipgloss.AdaptiveColor{Light: "#C92A2A", Dark: "#FF6B6B"}
	amber := lipgloss.AdaptiveColor{Light: "#B45309", Dark: "#F59E0B"}
	cyan := lipgloss.AdaptiveColor{Light: "#0E7490", Dark: "#22D3EE"}
	green := lipgloss.AdaptiveColor{Light: "#15803D", Dark: "#34D399"}
	gray := lipgloss.AdaptiveColor{Light: "#666666", Dark: "#9CA3AF"}
	subtle := lipgloss.AdaptiveColor{Light: "#999999", Dark: "#6B7280"}
	border := lipgloss.AdaptiveColor{Light: "#D4D4D8", Dark: "#3F3F46"}

	badge := func(c lipgloss.TerminalColor) lipgloss.Style {
		return r.NewStyle().
			Foreground(lipgloss.Color("#0F0F0F")).
			Background(c).
			Bold(true).
			Padding(0, 1)
	}

	return Theme{
		UseColor:     true,
		Brand:        r.NewStyle().Foreground(brand).Bold(true),
		BrandMark:    r.NewStyle().Foreground(cyan).Bold(true),
		Tagline:      r.NewStyle().Foreground(subtle).Italic(true),
		HeaderBorder: r.NewStyle().Foreground(border),

		CardBorder:   r.NewStyle().Foreground(border),
		BarFilled:    r.NewStyle().Foreground(green),
		BarEmpty:     r.NewStyle().Foreground(border),
		BarOverflow:  r.NewStyle().Foreground(red),

		SectionPrimary: r.NewStyle().Foreground(brand).Bold(true),
		SectionBonus:   r.NewStyle().Foreground(amber).Bold(true),
		SectionSubtle:  r.NewStyle().Foreground(subtle).Italic(true),

		SevHigh: badge(red),
		SevMed:  badge(amber),
		SevLow:  badge(cyan),
		SevInfo: badge(gray),

		ConfHigh: r.NewStyle().Foreground(green).Bold(true),
		ConfMed:  r.NewStyle().Foreground(amber).Bold(true),
		ConfLow:  r.NewStyle().Foreground(gray).Bold(true),

		Workload:   r.NewStyle().Foreground(brand).Bold(true),
		Title:      r.NewStyle().Bold(true),
		Detail:     r.NewStyle().Foreground(gray),
		Savings:    r.NewStyle().Foreground(green).Bold(true),
		BigSavings: r.NewStyle().Foreground(green).Bold(true),
		NoSavings:  r.NewStyle().Foreground(subtle),
		Muted:      r.NewStyle().Foreground(subtle),
		Divider:    r.NewStyle().Foreground(border),
		Disclosure: r.NewStyle().Foreground(amber),
		// The hyperlink already renders as clickable in modern terminals
		// (OSC 8); doubling that with `Underline(true)` makes lipgloss
		// emit per-character styling that confuses some terminals.
		CallToLink: r.NewStyle().Foreground(brand).Bold(true),
		OK:         r.NewStyle().Foreground(green).Bold(true),
	}
}

func plainTheme() Theme {
	plain := lipgloss.NewStyle()
	return Theme{
		UseColor:       false,
		Brand:          plain,
		BrandMark:      plain,
		Tagline:        plain,
		HeaderBorder:   plain,
		CardBorder:     plain,
		BarFilled:      plain,
		BarEmpty:       plain,
		BarOverflow:    plain,
		SectionPrimary: plain,
		SectionBonus:   plain,
		SectionSubtle:  plain,
		SevHigh:        plain,
		SevMed:         plain,
		SevLow:         plain,
		SevInfo:        plain,
		ConfHigh:       plain,
		ConfMed:        plain,
		ConfLow:        plain,
		Workload:       plain,
		Title:          plain,
		Detail:         plain,
		Savings:        plain,
		BigSavings:     plain,
		NoSavings:      plain,
		Muted:          plain,
		Divider:        plain,
		Disclosure:     plain,
		CallToLink:     plain,
		OK:             plain,
	}
}

// Hyperlink wraps url with an OSC 8 hyperlink escape so modern
// terminals (iTerm2, kitty, WezTerm, Ghostty, VSCode) render it as a
// clickable link. Falls back to plain text when colours are off.
func (t Theme) Hyperlink(label, url string) string {
	if !t.UseColor {
		return fmt.Sprintf("%s (%s)", label, url)
	}
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, label)
}

// DividerLine returns a horizontal rule the given width.
func (t Theme) DividerLine(width int) string {
	if width <= 0 {
		width = 64
	}
	return t.Divider.Render(repeat("─", width))
}

// SignalBar renders a horizontal ratio bar of the form
//
//	████████████░░░░░░░░  (have/want = 0.6)
//
// width is the total cell count; have and want are in the same unit
// and need not be normalized — the bar shows have/want as a fraction
// of width. When have > want (over-saturated) the overflow tail is
// drawn in BarOverflow so it visually screams.
//
// Returns a fixed-rune-width string regardless of color setting; the
// caller is responsible for any leading/trailing labels.
func (t Theme) SignalBar(have, want float64, width int) string {
	if width <= 0 {
		width = 20
	}
	if want <= 0 {
		// Degenerate: render an empty bar. Renderers should avoid
		// calling this when want is zero, but we tolerate it.
		return t.BarEmpty.Render(repeat("░", width))
	}
	ratio := have / want
	if ratio < 0 {
		ratio = 0
	}

	if ratio <= 1 {
		filled := int(ratio*float64(width) + 0.5)
		if filled > width {
			filled = width
		}
		return t.BarFilled.Render(repeat("█", filled)) +
			t.BarEmpty.Render(repeat("░", width-filled))
	}

	// Over-saturated: fill the whole bar in overflow tone so the eye
	// catches it. We don't try to encode the magnitude — the Note
	// (e.g. "10x burst") carries that.
	return t.BarOverflow.Render(repeat("█", width))
}

// SectionRule renders a labelled section divider using heavy hyphens
// so it visually outranks the regular dividers around the header and
// footer. Renders identically with or without color.
//
//	━━ <label> ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//
// padLeft/padRight let callers indent the rule to match surrounding
// content. accent picks the colour applied to the label and the rule.
func (t Theme) SectionRule(label string, width int, accent lipgloss.Style) string {
	if width <= 0 {
		width = 64
	}
	prefix := "━━ "
	rendered := accent.Render(prefix + label + " ")
	consumed := len([]rune(prefix + label + " "))
	remaining := width - consumed
	if remaining < 4 {
		remaining = 4
	}
	return rendered + accent.Render(repeat("━", remaining))
}

// SeverityBadge picks the right badge style for a severity string and
// renders the literal label.
func (t Theme) SeverityBadge(sev string) string {
	switch sev {
	case "HIGH":
		return t.SevHigh.Render(" HIGH ")
	case "MED":
		return t.SevMed.Render(" MED  ")
	case "LOW":
		return t.SevLow.Render(" LOW  ")
	default:
		return t.SevInfo.Render(" INFO ")
	}
}

// ConfidenceDots returns a fixed-width visual confidence indicator.
func (t Theme) ConfidenceDots(conf string) string {
	switch conf {
	case "high":
		return t.ConfHigh.Render("●●●") + " " + t.Muted.Render("high")
	case "medium":
		return t.ConfMed.Render("●●") + t.Muted.Render("○") + " " + t.Muted.Render("medium")
	case "low":
		return t.ConfLow.Render("●") + t.Muted.Render("○○") + " " + t.Muted.Render("low")
	default:
		return t.Muted.Render("○○○ unknown")
	}
}

// ConfidenceGlyph is a compact dot-only confidence indicator (no
// trailing word), used by the bonus-section one-liners where vertical
// density matters more than self-description.
func (t Theme) ConfidenceGlyph(conf string) string {
	switch conf {
	case "high":
		return t.ConfHigh.Render("●●●")
	case "medium":
		return t.ConfMed.Render("●●") + t.Muted.Render("○")
	case "low":
		return t.ConfLow.Render("●") + t.Muted.Render("○○")
	default:
		return t.Muted.Render("○○○")
	}
}

// IsTTY reports whether the given file is connected to a terminal.
// Centralised here so callers don't import isatty everywhere.
func IsTTY(f *os.File) bool {
	return isTTY(f)
}

func repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	out := make([]byte, 0, len(s)*n)
	for i := 0; i < n; i++ {
		out = append(out, s...)
	}
	return string(out)
}
