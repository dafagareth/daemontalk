package og

import (
	"bytes"
	"image/png"
	"os"
	"testing"
)

func TestRenderWithCover(t *testing.T) {
	var buf bytes.Buffer
	card := Card{
		Title:    "Ironi Pengumpulan Kode Sumber via Flashdisk dan Google Drive",
		Tags:     []string{"tools", "college", "git"},
		ReadTime: 5,
		Date:     "04 SEP 2026",
		Site:     "daemontalk.com",
		Cover:    "/static/images/posts/ironi-pengumpulan-kode-sumber-via-flashdisk-dan-google-drive-di-jurusan-teknologi/cover.png",
	}
	if err := Render(&buf, card); err != nil {
		t.Fatalf("Render with cover: %v", err)
	}

	rawBytes := buf.Bytes()
	_ = os.WriteFile("/home/dd/.gemini/antigravity-cli/brain/85d15c82-6877-4f62-9dec-5dc3ab6f9eb0/production_og_post.png", rawBytes, 0644)

	img, err := png.Decode(bytes.NewReader(rawBytes))
	if err != nil {
		t.Fatalf("decode PNG: %v", err)
	}
	b := img.Bounds()
	if b.Dx() != width || b.Dy() != height {
		t.Errorf("size: got %dx%d, want %dx%d", b.Dx(), b.Dy(), width, height)
	}
}

func TestRenderFallback(t *testing.T) {
	var buf bytes.Buffer
	card := Card{
		Title:    "DaemonTalk — Engineering & Systems Exploration",
		Tags:     []string{"SYSTEMS", "SOFTWARE", "LINUX"},
		ReadTime: 5,
		Date:     "04 SEP 2026",
		Site:     "daemontalk.com",
	}
	if err := Render(&buf, card); err != nil {
		t.Fatalf("Render fallback: %v", err)
	}

	rawBytes := buf.Bytes()
	_ = os.WriteFile("/home/dd/.gemini/antigravity-cli/brain/85d15c82-6877-4f62-9dec-5dc3ab6f9eb0/production_og_default.png", rawBytes, 0644)

	img, err := png.Decode(bytes.NewReader(rawBytes))
	if err != nil {
		t.Fatalf("decode PNG: %v", err)
	}
	b := img.Bounds()
	if b.Dx() != width || b.Dy() != height {
		t.Errorf("size: got %dx%d, want %dx%d", b.Dx(), b.Dy(), width, height)
	}
}

func TestRenderLongTitle(t *testing.T) {
	var buf bytes.Buffer
	long := "This is an extremely long blog post title that should wrap across multiple lines and eventually get truncated with an ellipsis when it exceeds the maximum allowed number of lines on the card"
	card := Card{Title: long, Site: "daemontalk.com"}
	if err := Render(&buf, card); err != nil {
		t.Fatalf("Render long title: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestWrap(t *testing.T) {
	face := newFace(pjsBoldFont, 48)
	defer face.Close()

	lines := wrap(face, "short", 1000)
	if len(lines) != 1 {
		t.Errorf("short text should be 1 line, got %d", len(lines))
	}

	lines = wrap(face, "one two three four five six seven eight nine ten eleven twelve", 50)
	if len(lines) > 4 {
		t.Errorf("should cap at 4 lines, got %d", len(lines))
	}
}

func TestSplitWords(t *testing.T) {
	got := splitWords("  hello   world\tfoo\n")
	want := []string{"hello", "world", "foo"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("word %d: got %q, want %q", i, got[i], want[i])
		}
	}
}
