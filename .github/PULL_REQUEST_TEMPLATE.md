### Summary of Changes
<!-- Brief summary of what this PR introduces (New dispatch, bug fix, typo correction, translation) -->

- **Type of Contribution**:
  - [ ] 📄 New Technical Dispatch (`content/posts/`)
  - [ ] ✏️ Article Correction / Code Snippet Fix
  - [ ] 🌐 Localization & Translation (`.id.md` / `.es.md`)
  - [ ] ⚡ Core Engine / Bug Fix (`Go`, `Templ`, `Tailwind`)

### Dispatch Checklist (For Article Submissions)
- [ ] Frontmatter YAML includes `title`, `slug`, `author`, `author_github`, `date`, `tags`, and `description`.
- [ ] If submitting a correction, added your GitHub handle to `contributors: ["your-handle"]`.
- [ ] Used standard Daemontalk markdown components (`> [!NOTE]`, ````stat````, ````tabs````, ````references````).
- [ ] Verified build and tests pass locally: `make build && go test -count=1 ./...`.
- [ ] No conversational AI filler; starts directly with the technical core.
