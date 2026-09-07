package post

import (
	"fmt"
	"html"
	"strings"
)

func renderCalloutHTML(kind, rawBody string) string {
	kind = strings.ToUpper(strings.TrimSpace(kind))
	if kind == "" {
		kind = "NOTE"
	}

	calloutClass := "post-callout-note"
	switch kind {
	case "TIP":
		calloutClass = "post-callout-tip"
	case "WARNING":
		calloutClass = "post-callout-warning"
	case "IMPORTANT":
		calloutClass = "post-callout-important"
	case "CAUTION":
		calloutClass = "post-callout-caution"
	default:
		kind = "NOTE"
		calloutClass = "post-callout-note"
	}

	bodyEsc := formatSimpleInline(html.EscapeString(strings.TrimSpace(rawBody)))
	bodyLines := strings.Split(bodyEsc, "\n")
	var formattedBody []string
	for _, l := range bodyLines {
		t := strings.TrimSpace(l)
		if t != "" {
			formattedBody = append(formattedBody, fmt.Sprintf("<p>%s</p>", t))
		}
	}

	return fmt.Sprintf(`
<div class="post-callout %s not-prose">
  <div class="callout-header">
    <span class="callout-label callout-title">%s</span>
  </div>
  <div class="callout-body">
    %s
  </div>
</div>
`, calloutClass, kind, strings.Join(formattedBody, "\n"))
}
