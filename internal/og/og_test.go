package og

import (
	"bytes"
	"image/png"
	"testing"
)

func TestRender(t *testing.T) {
	var buf bytes.Buffer
	card := Card{
		Title:    "Building stash: an encrypted secret vault in Go",
		Subtitle: "5 min read · go, cli, security",
		Site:     "daemontalk.com",
	}
	if err := Render(&buf, card); err != nil {
		t.Fatalf("Render: %v", err)
	}

	img, err := png.Decode(&buf)
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
	face := newFace(boldFont, 64)
	defer face.Close()

	lines := wrap(face, "short", 1000)
	if len(lines) != 1 {
		t.Errorf("short text should be 1 line, got %d", len(lines))
	}

	// Very narrow width forces many lines, capped at 5.
	lines = wrap(face, "one two three four five six seven eight nine ten", 50)
	if len(lines) > 5 {
		t.Errorf("should cap at 5 lines, got %d", len(lines))
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
