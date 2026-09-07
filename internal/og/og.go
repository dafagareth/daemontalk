package og

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed assets/PlusJakartaSans-Bold.ttf
var pjsBoldData []byte

//go:embed assets/Lora-Bold.ttf
var loraBoldData []byte

//go:embed assets/icon-light.png
var mascotIconData []byte

const (
	width  = 1200
	height = 630
)

var (
	cWhite      = color.RGBA{0xff, 0xff, 0xff, 0xff}
	cBorder     = color.RGBA{0x18, 0x18, 0x1b, 0xff}
	cBorderThin = color.RGBA{0xe4, 0xe4, 0xe7, 0xff}
	cText       = color.RGBA{0x09, 0x09, 0x0b, 0xff}
	cSubtle     = color.RGBA{0xa1, 0xa1, 0xaa, 0xff}

	pjsBoldFont  *opentype.Font
	loraBoldFont *opentype.Font
	mascotIcon   image.Image
)

func init() {
	var err error
	pjsBoldFont, err = opentype.Parse(pjsBoldData)
	if err != nil {
		panic("og: parse Plus Jakarta Sans font: " + err.Error())
	}
	loraBoldFont, err = opentype.Parse(loraBoldData)
	if err != nil {
		panic("og: parse Lora font: " + err.Error())
	}
	mascotIcon, _, err = image.Decode(bytes.NewReader(mascotIconData))
	if err != nil {
		panic("og: decode mascot icon: " + err.Error())
	}
}

type Card struct {
	Title    string
	Subtitle string
	Tags     []string
	ReadTime int
	Date     string
	Site     string
	Cover    string
}

func Render(w io.Writer, c Card) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	drawRect(img, 0, 0, width, height, cWhite)

	bx, by, bw, bh := 36, 36, 1128, 558
	for x := bx; x < bx+bw; x++ {
		img.Set(x, by, cBorder)
		img.Set(x, by+bh-1, cBorder)
	}
	for y := by; y < by+bh; y++ {
		img.Set(bx, y, cBorder)
		img.Set(bx+bw-1, y, cBorder)
	}

	headY := by + 56
	for x := bx; x < bx+bw; x++ {
		img.Set(x, headY, cBorderThin)
	}

	logoW, logoH := 42, 31
	logoDst := image.Rect(bx+24, by+13, bx+24+logoW, by+13+logoH)
	xdraw.BiLinear.Scale(img, logoDst, mascotIcon, mascotIcon.Bounds(), xdraw.Over, nil)

	loraLogoFace := newFace(loraBoldFont, 24)
	textStartX := bx + 76
	drawText(img, loraLogoFace, "daemontalk", textStartX, by+36, cText)
	loraLogoFace.Close()

	cellH := 84
	cellY := by + bh - cellH
	for x := bx; x < bx+bw; x++ {
		img.Set(x, cellY, cBorder)
	}

	colW := bw / 4
	for col := 1; col < 4; col++ {
		colX := bx + col*colW
		for y := cellY; y < cellY+cellH; y++ {
			img.Set(colX, y, cBorderThin)
		}
	}

	labelFace := newFace(pjsBoldFont, 13)
	valFace := newFace(pjsBoldFont, 19)

	tagsVal := formatTags(c.Tags, c.Subtitle)
	readTimeVal := "5 MINUTES"
	if c.ReadTime > 0 {
		readTimeVal = fmt.Sprintf("%d MINUTES", c.ReadTime)
	}
	dateVal := strings.ToUpper(c.Date)
	if dateVal == "" {
		dateVal = strings.ToUpper(time.Now().Format("02 JAN 2006"))
	}
	siteVal := c.Site
	if siteVal == "" {
		siteVal = "daemontalk.com"
	}

	drawText(img, labelFace, "DOMAIN", bx+24, cellY+30, cSubtle)
	drawText(img, valFace, tagsVal, bx+24, cellY+60, cText)

	drawText(img, labelFace, "READ TIME", bx+colW+24, cellY+30, cSubtle)
	drawText(img, valFace, readTimeVal, bx+colW+24, cellY+60, cText)

	drawText(img, labelFace, "DATE", bx+colW*2+24, cellY+30, cSubtle)
	drawText(img, valFace, dateVal, bx+colW*2+24, cellY+60, cText)

	drawText(img, labelFace, "SOURCE", bx+colW*3+24, cellY+30, cSubtle)
	drawText(img, valFace, siteVal, bx+colW*3+24, cellY+60, cText)

	labelFace.Close()
	valFace.Close()

	coverImg := loadCover(c.Cover)
	if coverImg != nil {

		coverDst := image.Rect(bx+1, headY+1, bx+bw-1, cellY)
		drawCoverFit(img, coverDst, coverImg)

		gradStartY := headY + 140
		drawBottomGradient(img, bx+1, gradStartY, bx+bw-1, cellY, 210)

		titleFace := newFace(pjsBoldFont, 48)
		maxW := bw - 120
		lines := wrap(titleFace, c.Title, maxW)
		lineH := 62
		totalTextH := (len(lines)-1)*lineH + 46
		textStartY := cellY - 36 - totalTextH + 38

		for i, ln := range lines {

			drawText(img, titleFace, ln, bx+44+1, textStartY+i*lineH+2, color.RGBA{0, 0, 0, 160})
			drawText(img, titleFace, ln, bx+44, textStartY+i*lineH, cWhite)
		}
		titleFace.Close()
	} else {

		titleFace := newFace(pjsBoldFont, 52)
		maxW := bw - 88
		lines := wrap(titleFace, c.Title, maxW)
		lineH := 70
		availableH := cellY - headY
		blockH := len(lines) * lineH
		startY := headY + (availableH-blockH)/2 + 50
		if startY < headY+60 {
			startY = headY + 60
		}

		for i, ln := range lines {
			drawText(img, titleFace, ln, bx+44, startY+i*lineH, cText)
		}
		titleFace.Close()
	}

	return png.Encode(w, img)
}

func formatTags(tags []string, fallback string) string {
	if len(tags) > 0 {
		var parts []string
		for i, t := range tags {
			if i >= 2 {
				break
			}
			parts = append(parts, "#"+strings.ToUpper(strings.TrimSpace(t)))
		}
		return strings.Join(parts, ", ")
	}
	if fallback != "" {
		parts := strings.Split(fallback, "·")
		if len(parts) > 1 {
			raw := strings.TrimSpace(parts[1])
			items := strings.Split(raw, ",")
			var formatted []string
			for i, it := range items {
				if i >= 2 {
					break
				}
				formatted = append(formatted, "#"+strings.ToUpper(strings.TrimSpace(it)))
			}
			if len(formatted) > 0 {
				return strings.Join(formatted, ", ")
			}
		}
	}
	return "#SYSTEMS, #TECH"
}

func loadCover(cover string) image.Image {
	if cover == "" {
		return nil
	}

	if strings.HasPrefix(cover, "http://") || strings.HasPrefix(cover, "https://") {
		client := &http.Client{Timeout: 4 * time.Second}
		resp, err := client.Get(cover)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			img, _, err := image.Decode(resp.Body)
			if err == nil {
				return img
			}
		}
		return nil
	}

	clean := filepath.Clean(cover)
	trimmed := strings.TrimPrefix(clean, "/")
	trimmedStatic := strings.TrimPrefix(trimmed, "static/")

	candidates := []string{
		clean,
		trimmed,
		filepath.Join("web", trimmed),
		filepath.Join("web", "static", trimmedStatic),
		filepath.Join("..", "web", trimmed),
		filepath.Join("..", "..", "web", trimmed),
		filepath.Join("..", "..", "web", "static", trimmedStatic),
	}
	for _, path := range candidates {
		f, err := os.Open(path)
		if err == nil {
			img, _, err := image.Decode(f)
			_ = f.Close()
			if err == nil {
				return img
			}
		}
	}

	return nil
}

func drawCoverFit(dst xdraw.Image, dstRect image.Rectangle, src image.Image) {
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	dstW := dstRect.Dx()
	dstH := dstRect.Dy()

	srcAspect := float64(srcW) / float64(srcH)
	dstAspect := float64(dstW) / float64(dstH)

	var cropRect image.Rectangle
	if srcAspect > dstAspect {
		cropW := int(float64(srcH) * dstAspect)
		cropH := srcH
		cropX := srcBounds.Min.X + (srcW-cropW)/2
		cropY := srcBounds.Min.Y
		cropRect = image.Rect(cropX, cropY, cropX+cropW, cropY+cropH)
	} else {
		cropW := srcW
		cropH := int(float64(srcW) / dstAspect)
		cropX := srcBounds.Min.X
		cropY := srcBounds.Min.Y + (srcH-cropH)/2
		cropRect = image.Rect(cropX, cropY, cropX+cropW, cropY+cropH)
	}

	xdraw.BiLinear.Scale(dst, dstRect, src, cropRect, xdraw.Over, nil)
}

func drawBottomGradient(img *image.RGBA, x0, y0, x1, y1 int, maxAlpha uint8) {
	h := float64(y1 - y0)
	for y := y0; y < y1; y++ {
		if y < 0 || y >= img.Bounds().Dy() {
			continue
		}
		t := float64(y-y0) / h
		alpha := int(float64(maxAlpha) * (t*t*0.9 + t*0.1))
		invA := 255 - alpha
		for x := x0; x < x1; x++ {
			if x < 0 || x >= img.Bounds().Dx() {
				continue
			}
			c := img.RGBAAt(x, y)
			nr := uint8(int(c.R) * invA / 255)
			ng := uint8(int(c.G) * invA / 255)
			nb := uint8(int(c.B) * invA / 255)
			img.SetRGBA(x, y, color.RGBA{nr, ng, nb, 255})
		}
	}
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
	if len(lines) > 4 {
		lines = lines[:4]
		lines[3] = truncate(face, lines[3]+"…", maxW)
	}
	return lines
}

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
