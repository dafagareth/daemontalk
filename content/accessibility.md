# Accessibility Statement

**Last updated: August 2026**

Daemontalk is built to be accessible, performant, and legible across all devices, screen readers, and terminal environments, adhering to the Web Content Accessibility Guidelines (**WCAG 2.1 Level AA**).

## Keyboard Navigation & Hotkeys

All interactive elements support keyboard navigation with visible focus indicators:

- `?` : Open the Keyboard Shortcuts Cheat Sheet.
- `/` : Instant search bar focus.
- `t` : Switch color themes (Light, Dark, Sepia).
- `j` / `k` : Navigate previous and next articles.
- `Esc` : Dismiss active modals and search popups.
- `Tab` : Step sequentially through all links and actions with high-contrast focus rings.

## Typography & Contrast

- **Readable Typefaces**: Body text is set in **Source Serif 4** for optimal reading rhythm, user interface elements in **Plus Jakarta Sans**, and code snippets in **JetBrains Mono**.
- **Font Scaling**: Articles include on-page text scaling controls (`A-` and `A+`) that reflow text vertically without clipping or breaking layout.
- **High Contrast Ratios**: Text-to-background contrast exceeds 4.5:1 in both light and dark modes.
- **Warm Screen Tint**: An integrated warm tint slider reduces blue light exposure for readers with light sensitivity.

## Screen Readers & Semantic Structure

- The site uses standard HTML5 semantic elements (`<main>`, `<nav>`, `<article>`, `<header>`, `<footer>`, `<aside>`).
- Interactive icons and decorative graphics include explicit ARIA attributes (`aria-label`, `aria-hidden`).
- Architecture diagrams and technical figures include descriptive alternative text (`alt`).
- Respects system `prefers-reduced-motion` preferences by disabling non-essential transition animations.

## Terminal & Headless Access

For users working in text-only environments, console screens, or braille terminals, the entire site archive and dispatches can be accessed without a web browser via SSH:

```bash
$ ssh ssh.daemontalk.com -p 2222
```

## Feedback & Assistance

If you encounter any accessibility barriers on this site, please email **realdaemontalk@gmail.com**. We treat accessibility bug reports with high priority.
