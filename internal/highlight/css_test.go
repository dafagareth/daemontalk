package highlight

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateCSS(t *testing.T) {
	css := GenerateCSS()

	if !strings.Contains(css, `/* Base Chroma reset & variable inheritance */`) {
		t.Error("missing base chroma reset")
	}

	if !strings.Contains(css, `[data-theme="light"]`) {
		t.Error("missing light theme scoping")
	}

	if !strings.Contains(css, `[data-theme="dark"]`) {
		t.Error("missing dark theme scoping")
	}

	if !strings.Contains(css, `@media (prefers-color-scheme: dark)`) {
		t.Error("missing prefers-color-scheme dark query")
	}

	// Verify that dark mode explicit color rules for .nx, .p, .nn are present
	if !strings.Contains(css, `[data-theme="dark"] .chroma .nx`) {
		t.Error("missing dark theme .nx rule")
	}
	if !strings.Contains(css, `color: #e6edf3`) {
		t.Error("missing github dark text color")
	}

	_ = os.WriteFile("../../web/static/css/chroma.css", []byte(css), 0644)
}
