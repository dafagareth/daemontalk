# Accessibility Statement

**Effective Date: September 12, 2026 · Standard: WCAG 2.1 Level AA**

Daemontalk is committed to making technical knowledge universally accessible. We target **WCAG 2.1 Level AA** conformance and design the platform to work well for users navigating via keyboard, screen readers, and assistive technologies.

---

## Keyboard Navigation

The entire platform is operable without a mouse. Tab and Shift+Tab step through interactive elements in logical DOM order, and focus indicators are visible at all times. Pressing `/` focuses the global search from anywhere on the page, `?` opens the keyboard shortcut reference, `j` and `k` navigate between dispatches, `t` cycles the theme, and `Esc` closes any open dialog or overlay. No keyboard traps exist.

---

## Typography & Contrast

Body text meets a minimum 4.5:1 contrast ratio; large headings meet 3:1. Long-form prose is set in Lora (serif), interface controls in Plus Jakarta Sans, and code blocks in JetBrains Mono. Readers can adjust prose size across three preset steps via the reading controls. The platform supports browser zoom up to 200% without horizontal overflow or loss of functionality.

---

## Semantic Markup

Page templates use semantic HTML5 landmarks (`header`, `nav`, `main`, `article`, `aside`, `footer`) to allow screen readers to jump directly between sections. Heading hierarchy is maintained without skipped levels. Interactive components carry appropriate WAI-ARIA attributes (`aria-expanded`, `aria-controls`, `aria-hidden`, `aria-live`). Decorative elements are marked `aria-hidden="true"` and meaningful images carry descriptive `alt` text.

---

## Motion & Cognitive Comfort

The platform respects the `prefers-reduced-motion` OS setting and suppresses transitions and animated transforms accordingly. No content flashes more than three times per second. There are no auto-playing media, countdown timers, or unsolicited overlays.

---

## Known Limitations

Some deep-dive kernel articles contain raw terminal output, ANSI streams, or hex dumps for which full accessible alternatives are not always practical. In those cases we provide plain-text descriptions alongside the raw blocks where possible.

---

## Feedback

If you encounter an accessibility barrier, please email **realdaemontalk@gmail.com** with the subject line `[Accessibility Barrier]`, the affected URL, your assistive technology and browser, and a brief description. Daemontalk is maintained by a single operator and we will address reported barriers as promptly as possible.
