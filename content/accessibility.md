# Accessibility Statement

**Last updated: September 4, 2026** · **Standard: WCAG 2.1 Level AA**

Daemontalk is committed to providing an inclusive, fast, and legible reading and browsing experience for all individuals, including users with visual, motor, auditory, or cognitive disabilities, as well as those operating in headless or screen reader environments.

---

## Universal Design Principles

This platform is designed under universal accessibility principles aligned with the international Web Content Accessibility Guidelines (**WCAG 2.1 Level AA**). We ensure that navigational ease, typographic clarity, and visual contrast remain paramount without sacrificing performance.

## Comprehensive Keyboard Navigation

All platform features and interactive components support full keyboard navigation without requiring a pointing device:

**Global Shortcuts Guide (`?`)**: Opens an interactive modal cheat sheet detailing all available keyboard hotkeys.

**Instant Search (`/`)**: Focuses immediately on the site-wide search input.

**Theme Toggle (`t`)**: Cycles visual color schemes between Light and Dark modes.

**Sequential Post Navigation (`j` and `k`)**: Steps forward and backward through published dispatches.

**Modal Dismissal (`Esc`)**: Closes any active search popup, shortcut modal, or drawer menu.

**Linear Tab Stepping (`Tab` & `Shift+Tab`)**: Navigates sequentially through links, buttons, and form inputs with high-contrast focus rings for clear visual tracking.

## Ergonomic Typography & Visual Contrast

**High-Legibility Font Pairing**: Body prose is rendered in **Source Serif 4** with generous line height for sustained technical reading, interface controls in **Plus Jakarta Sans**, and code blocks in **JetBrains Mono**.

**Dynamic Font Scaling**: Every dispatch includes on-page text scaling controls (`A-` and `A+`) that reflow prose vertically without clipping words or breaking responsive layouts.

**Calibrated Contrast Ratios**: All text-to-background color combinations maintain a contrast ratio exceeding 4.5:1 for normal text and 3:1 for large display headers in both light and dark themes.

**Warm Screen Tint**: An integrated warm tint slider filters high-frequency blue light to alleviate eye strain during low-light reading sessions.

## Semantic Structure & Screen Reader Compatibility

**HTML5 Landmark Semantics**: Document layouts use strict HTML5 structural tags (`<main>`, `<nav>`, `<article>`, `<header>`, `<footer>`, `<aside>`, `<section>`) allowing assistive screen readers (NVDA, VoiceOver, JAWS, Orca) to navigate landmarks effortlessly.

**Explicit ARIA Attributes**: Interactive controls, iconography, asynchronous loading states, and dialog windows are equipped with accurate ARIA labels and states (`aria-label`, `aria-expanded`, `aria-hidden`, `role="dialog"`).

**Descriptive Alternative Text**: System blueprints, performance charts, and flow diagrams include contextual `alt` text explaining their technical meaning.

## Motion Reduction & Cognitive Comfort

**Reduced Motion Support**: The site automatically detects system `prefers-reduced-motion` settings and eliminates non-essential transition animations for readers with vestibular sensitivities.

**Distraction-Free Environment**: The platform contains zero rapidly flashing elements, zero autoplaying media, and zero invasive marketing popups.

## Headless & Terminal Accessibility

For users operating in text-only environments, console displays, or braille terminals, the entire technical archive is accessible directly over SSH without requiring a graphical web browser:

```bash
$ ssh ssh.daemontalk.com -p 2222
```

## Feedback & Accessibility Assistance

We continuously audit and improve the accessibility of this platform. If you encounter any barriers, difficult-to-read text, or keyboard traps, please reach out directly via email at: **realdaemontalk@gmail.com**. All accessibility inquiries are treated with high priority.
