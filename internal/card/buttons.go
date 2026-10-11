package card

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type btnTheme struct{ bg, border, text string }

var btnThemes = map[string]btnTheme{
	"dark":  {bg: "#21262d", border: "#3d444d", text: "#c9d1d9"},
	"light": {bg: "#f6f8fa", border: "#d1d9e0", text: "#24292f"},
}

const btnFont = "-apple-system,BlinkMacSystemFont,'Segoe UI',Helvetica,Arial,sans-serif"

const (
	btnH    = 32.0
	btnFS   = 14.0
	btnPadL = 12.0
	btnIC   = 16.0
	btnGap  = 8.0
	btnPadR = 12.0
	btnRX   = ghRadius
)

// glyphAdv and glyphRSB are per-glyph metrics for SF Pro Text Medium — the font
// that -apple-system resolves to at the button's 14px font-size — measured in
// 1000-em units via CoreText. (The previous Helvetica AFM table underestimated
// SF Pro's width by ~6px on short labels, leaving the buttons with noticeably
// more left than right padding.) glyphAdv is the advance width (pen travel);
// glyphRSB is the trailing right side bearing (advance minus ink right edge),
// used to center the visible ink rather than the advance box. Unmapped runes use
// the *Default fallbacks.
const (
	glyphAdvDefault = 600.0
	glyphRSBDefault = 40.0
)

var glyphAdv = map[rune]float64{
	' ': 262, '!': 317, '"': 500, '#': 634, '$': 634, '%': 955, '&': 713, '\'': 305,
	'(': 388, ')': 388, '*': 466, '+': 634, ',': 305, '-': 466, '.': 305, '/': 305,
	'0': 639, '1': 471, '2': 607, '3': 632, '4': 649, '5': 624, '6': 643, '7': 571,
	'8': 648, '9': 643, ':': 305, ';': 305, '<': 634, '=': 634, '>': 634, '?': 518,
	'@': 911, 'A': 685, 'B': 660, 'C': 714, 'D': 723, 'E': 595, 'F': 571, 'G': 741,
	'H': 746, 'I': 275, 'J': 550, 'K': 666, 'L': 567, 'M': 875, 'N': 741, 'O': 766,
	'P': 639, 'Q': 766, 'R': 658, 'S': 640, 'T': 633, 'U': 735, 'V': 680, 'W': 972,
	'X': 687, 'Y': 664, 'Z': 656, '[': 388, '\\': 305, ']': 388, '^': 634, '_': 589,
	'`': 489, 'a': 555, 'b': 617, 'c': 559, 'd': 617, 'e': 571, 'f': 368, 'g': 612,
	'h': 594, 'i': 251, 'j': 251, 'k': 554, 'l': 258, 'm': 882, 'n': 589, 'o': 591,
	'p': 613, 'q': 613, 'r': 390, 's': 528, 't': 371, 'u': 589, 'v': 546, 'w': 790,
	'x': 534, 'y': 552, 'z': 538, '{': 388, '|': 263, '}': 388, '~': 634,
}

var glyphRSB = map[rune]float64{
	' ': 262, '!': 78, '"': 94, '#': 4, '$': 60, '%': 54, '&': 2, '\'': 93,
	'(': 39, ')': 104, '*': 53, '+': 47, ',': 91, '-': 64, '.': 75, '/': -11,
	'0': 54, '1': 112, '2': 56, '3': 52, '4': 46, '5': 50, '6': 51, '7': 44,
	'8': 51, '9': 51, ':': 75, ';': 75, '<': 59, '=': 53, '>': 59, '?': 41,
	'@': 44, 'A': 23, 'B': 41, 'C': 44, 'D': 49, 'E': 58, 'F': 42, 'G': 51,
	'H': 77, 'I': 77, 'J': 77, 'K': 16, 'L': 39, 'M': 77, 'N': 77, 'O': 48,
	'P': 41, 'Q': 48, 'R': 35, 'S': 44, 'T': 35, 'U': 76, 'V': 23, 'W': 32,
	'X': 33, 'Y': 23, 'Z': 58, '[': 27, '\\': -11, ']': 117, '^': 45, '_': -11,
	'`': 140, 'a': 62, 'b': 39, 'c': 32, 'd': 70, 'e': 39, 'f': 27, 'g': 65,
	'h': 60, 'i': 53, 'j': 53, 'k': 9, 'l': 70, 'm': 60, 'n': 60, 'o': 39,
	'p': 40, 'q': 64, 'r': 14, 's': 39, 't': 37, 'u': 65, 'v': 19, 'w': 21,
	'x': 23, 'y': 20, 'z': 53, '{': 27, '|': 77, '}': 27, '~': 53,
}

// iconLSB is each icon's left ink bearing in px within its 16px box: the
// whitespace before the glyph's first drawn pixel. Used to center the icon ink
// rather than its layout box. Derived from the icon path geometry.
var iconLSB = map[string]float64{
	"github":   1.0,
	"linkedin": 0.667,
	"email":    0.667,
}

func labelWidth(label string) float64 {
	var units float64
	for _, r := range label {
		if a, ok := glyphAdv[r]; ok {
			units += a
		} else {
			units += glyphAdvDefault
		}
	}
	return units / 1000 * btnFS
}

// labelLastRSB returns the right side bearing (px at btnFS) of the label's
// final glyph — the trailing whitespace between the text's advance box and its
// visible ink.
func labelLastRSB(label string) float64 {
	r := []rune(label)
	if len(r) == 0 {
		return 0
	}
	last := r[len(r)-1]
	if b, ok := glyphRSB[last]; ok {
		return b / 1000 * btnFS
	}
	return glyphRSBDefault / 1000 * btnFS
}

type Button struct{ Label, Icon string }

func lucideGitHub(x, y, size float64, color string) string {
	s := size / 24
	return fmt.Sprintf(`<g transform="translate(%.2f %.2f) scale(%.4f)" fill="none" stroke="%s" `+
		`stroke-width="2" stroke-linecap="round" stroke-linejoin="round">`+
		`<path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"/>`+
		`<path d="M9 18c-4.51 2-5-2-7-2"/></g>`, x, y, s, color)
}

func lucideLinkedIn(x, y, size float64, color string) string {
	s := size / 24
	return fmt.Sprintf(`<g transform="translate(%.2f %.2f) scale(%.4f)" fill="none" stroke="%s" `+
		`stroke-width="2" stroke-linecap="round" stroke-linejoin="round">`+
		`<path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"/>`+
		`<rect width="4" height="12" x="2" y="9"/><circle cx="4" cy="4" r="2"/></g>`, x, y, s, color)
}

func lucideMail(x, y, size float64, color string) string {
	s := size / 24
	return fmt.Sprintf(`<g transform="translate(%.2f %.2f) scale(%.4f)" fill="none" stroke="%s" `+
		`stroke-width="2" stroke-linecap="round" stroke-linejoin="round">`+
		`<path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"/>`+
		`<rect x="2" y="4" width="20" height="16" rx="2"/></g>`, x, y, s, color)
}

func btnNaturalWidth(label string) float64 {
	return math.Ceil(btnPadL + btnIC + btnGap + labelWidth(label) + btnPadR)
}

func renderButton(themeName, label, icon string, width float64) string {
	t := btnThemes[themeName]
	// Center the visible ink (icon's left edge ↔ text's right edge) rather than
	// the advance box, so horizontal padding stays symmetric for any label/icon.
	adv := labelWidth(label)
	lastRSB := labelLastRSB(label)
	iL := iconLSB[icon]
	iconX := (width - btnIC - btnGap - adv + lastRSB - iL) / 2
	iconY := (btnH - btnIC) / 2
	textX := iconX + btnIC + btnGap
	textY := btnH/2 + btnFS*0.34

	var mark string
	switch icon {
	case "github":
		mark = lucideGitHub(iconX, iconY, btnIC, t.text)
	case "linkedin":
		mark = lucideLinkedIn(iconX, iconY, btnIC, t.text)
	case "email":
		mark = lucideMail(iconX, iconY, btnIC, t.text)
	}

	var o strings.Builder
	o.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" `+
		`viewBox="0 0 %.0f %.0f" font-family="%s" role="img" aria-label="%s">`,
		width, btnH, width, btnH, btnFont, esc(label)))
	o.WriteString(fmt.Sprintf(`<rect x="0.5" y="0.5" width="%.0f" height="%.0f" rx="%g" fill="%s" stroke="%s"/>`,
		width-1, btnH-1, btnRX, t.bg, t.border))
	o.WriteString(mark)
	o.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%g" font-weight="500" fill="%s">%s</text>`,
		textX, textY, btnFS, t.text, esc(label)))
	o.WriteString(`</svg>`)
	return o.String()
}

// WriteButtons renders the connect buttons for both themes into
// dir/btn_<file>_<theme>.svg.
func WriteButtons(dir string) error {
	btns := []Button{
		{Label: "GitHub", Icon: "github"},
		{Label: "LinkedIn", Icon: "linkedin"},
		{Label: "Email", Icon: "email"},
	}
	var width float64
	for _, b := range btns {
		if w := btnNaturalWidth(b.Label); w > width {
			width = w
		}
	}
	for _, b := range btns {
		for _, themeName := range []string{"dark", "light"} {
			path := fmt.Sprintf("%s/btn_%s_%s.svg", dir, b.Icon, themeName)
			if err := os.WriteFile(path, []byte(renderButton(themeName, b.Label, b.Icon, width)+"\n"), 0o644); err != nil {
				return fmt.Errorf("write %s: %w", path, err)
			}
		}
	}
	return nil
}
