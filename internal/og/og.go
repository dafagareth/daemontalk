// Package og renders Open Graph share images (1200×630 PNG) for blog posts
// that don't ship their own cover image. Cards are drawn entirely in Go using
// the embedded Go fonts, so no external assets or headless browser are needed.
package og

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	width  = 1200
	height = 630
	padX   = 90
)

var (
	bgTop   = color.RGBA{0x0f, 0x17, 0x2a, 0xff} // slate-900
	bgBot   = color.RGBA{0x1a, 0x1f, 0x35, 0xff}
	accent  = color.RGBA{0x3b, 0x82, 0xf6, 0xff} // blue-500
	fgTitle = color.RGBA{0xf1, 0xf5, 0xf9, 0xff} // slate-100
	fgMuted = color.RGBA{0x64, 0x74, 0x8b, 0xff} // slate-500
	fgLink  = color.RGBA{0x60, 0xa5, 0xfa, 0xff} // blue-400

	boldFont    = mustParse(gobold.TTF)
	regularFont = mustParse(goregular.TTF)
)

func mustParse(ttf []byte) *opentype.Font {
	f, err := opentype.Parse(ttf)
	if err != nil {
		panic("og: parse font: " + err.Error())
	}
	return f
}

// Card holds the text rendered onto the share image.
type Card struct {
	Title    string
	Subtitle string // e.g. "5 min read · Go, CLI"
	Site     string // e.g. "dafagareth.dev"
}

// Render writes the card as a PNG to w.
func Render(w io.Writer, c Card) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	drawVerticalGradient(img, bgTop, bgBot)

	// Left accent bar
	drawRect(img, padX-30, 110, 5, height-220, accent)

	titleFace := newFace(boldFont, 64)
	defer titleFace.Close()
	maxW := width - padX - 100

	// Wrap the title and vertically center the block.
	lines := wrap(titleFace, c.Title, maxW)
	lineH := 84
	blockH := len(lines) * lineH
	startY := (height-blockH)/2 + 56
	if startY < 160 {
		startY = 160
	}
	for i, ln := range lines {
		drawText(img, titleFace, ln, padX, startY+i*lineH, fgTitle)
	}

	// Subtitle under the title block
	if c.Subtitle != "" {
		subFace := newFace(regularFont, 30)
		defer subFace.Close()
		drawText(img, subFace, c.Subtitle, padX, startY+len(lines)*lineH+30, fgMuted)
	}

	// Site label bottom-left
	if c.Site != "" {
		siteFace := newFace(boldFont, 32)
		defer siteFace.Close()
		drawText(img, siteFace, c.Site, padX, height-60, fgLink)
	}

	return png.Encode(w, img)
}

func newFace(f *opentype.Font, size float64) font.Face {
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic("og: new face: " + err.Error())
	}
	return face
}

// wrap greedily splits text into lines that fit within maxW pixels.
func wrap(face font.Face, text string, maxW int) []string {
	words := splitWords(text)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	line := words[0]
	for _, word := range words[1:] {
		test := line + " " + word
		if textWidth(face, test) > maxW {
			lines = append(lines, line)
			line = word
		} else {
			line = test
		}
	}
	lines = append(lines, line)
	// Cap at 5 lines; truncate the last with an ellipsis if needed.
	if len(lines) > 5 {
		lines = lines[:5]
		lines[4] = truncate(face, lines[4]+"…", maxW)
	}
	return lines
}

// truncate drops trailing runes until the string (with a trailing ellipsis)
// fits within maxW. Operates on runes so multi-byte characters aren't split.
func truncate(face font.Face, s string, maxW int) string {
	runes := []rune(strings.TrimSuffix(s, "…"))
	for len(runes) > 0 {
		candidate := string(runes) + "…"
		if textWidth(face, candidate) <= maxW {
			return candidate
		}
		runes = runes[:len(runes)-1]
	}
	return "…"
}

func splitWords(s string) []string {
	var words []string
	start := -1
	for i, r := range s {
		if r == ' ' || r == '\t' || r == '\n' {
			if start >= 0 {
				words = append(words, s[start:i])
				start = -1
			}
		} else if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		words = append(words, s[start:])
	}
	return words
}

func textWidth(face font.Face, s string) int {
	return font.MeasureString(face, s).Round()
}

func drawText(img *image.RGBA, face font.Face, s string, x, y int, col color.Color) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(s)
}

func drawRect(img *image.RGBA, x, y, w, h int, col color.Color) {
	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			img.Set(xx, yy, col)
		}
	}
}

func drawVerticalGradient(img *image.RGBA, top, bot color.RGBA) {
	for y := 0; y < height; y++ {
		t := float64(y) / float64(height-1)
		r := uint8(float64(top.R)*(1-t) + float64(bot.R)*t)
		g := uint8(float64(top.G)*(1-t) + float64(bot.G)*t)
		b := uint8(float64(top.B)*(1-t) + float64(bot.B)*t)
		row := color.RGBA{r, g, b, 0xff}
		for x := 0; x < width; x++ {
			img.Set(x, y, row)
		}
	}
}
