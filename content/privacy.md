# Privacy Policy

**Last updated: August 30, 2026** · **Version: 2.1**

Daemontalk operates under a strict zero-tracker philosophy. We believe technical publications and developer platforms must prioritize reader sovereignty, architectural transparency, and computational minimalism. We do not track, profile, monetize, or harvest personal data.

---

## 1. Zero Tracking & Self-Hosted Infrastructure

**No Behavioral Trackers**: This site operates without Google Analytics, Meta/Facebook Pixels, tracking beacons, advertising cookies, or browser fingerprinting scripts.

**No Third-Party Ad Networks**: We do not display commercial advertisements, sponsor tracking banners, or affiliate trackers, and we never sell user data to data brokers.

**Self-Hosted Assets**: All stylesheets (CSS), scripts (JavaScript), fonts, and media are served directly from our dedicated Go binary server, eliminating third-party surveillance vectors.

## 2. Legal Bases for Data Processing (GDPR & PDP)

Under international data protection regulations (including the EU General Data Protection Regulation and Indonesian Law No. 27/2022 on Personal Data Protection), we process minimal data strictly under the following legal bases:

**Consent**: When you explicitly choose to authenticate via GitHub OAuth, post a technical topic, submit an article comment, or send a message via the contact form.

**Legitimate Interest**: To maintain server uptime, mitigate DDoS and brute-force attacks via rate limiting, and debug system infrastructure using ephemeral, rotating connection logs with a strict 14-day retention window.

**Contractual & Operational Necessity**: To maintain your active login session and associate your verified author identity with your submitted forum contributions.

## 3. GitHub OAuth Authentication & Minimal Data Profile

Daemontalk provides optional member sign-in powered by GitHub OAuth to enable verified community identity across discussions, topics, and technical comments.

**Minimal Data Collected**: When you sign in via GitHub, we request only the minimal necessary read permissions (`read:user`, `user:email`). We store strictly your unique numerical GitHub User ID, public username handle (e.g. `octocat`), public display name, public avatar URL, public GitHub profile URL, and primary verified email address (used strictly for session identity validation and never shared or sent marketing emails).

**What We Never Access**: We never request, access, or store your private repositories, source code, SSH keys, billing details, organization secrets, or write permissions to your GitHub account.

**Session Tokens**: Login sessions are managed via cryptographically hashed random tokens (SHA-256) stored in standard HTTP-only, `SameSite=Lax` cookies. We do not use third-party session tracking trackers.

## 4. Client-Side Storage & Local Preferences

We utilize standard client-side browser storage mechanisms (`localStorage` and `sessionStorage`) exclusively on your local device to preserve your personal viewing preferences:

**Saved Dispatches (`bookmarks`)**: The list of bookmarked posts stored in your local reading ledger.

**Theme Selection (`theme`)**: Your active color scheme preference (*Light or Dark mode*).

**Typography & Accessibility**: Font size scaling, font family selection (*Serif or Sans*), and warm screen tint intensity.

**Visited Dispatches Indicator (`readPosts`)**: A local array of recently visited article slugs (capped at 200 items) used solely to indicate previously read content on index pages.

**Session Animations (`visited`)**: A temporary session token in `sessionStorage` to prevent redundant entry animations on subsequent page navigations.

**Data Sovereignty**: None of this client-side preference data is ever synced to our backend database or any third-party entity without your explicit action. You may purge this data at any time by clearing your browser's site storage.

## 5. Ephemeral Server Logs & Telemetry

When connecting to Daemontalk over HTTP or SSH, our self-hosted server records standard connection metadata including connecting IP address, request path, HTTP method, timestamp, User-Agent string, Referrer header (if transmitted by your client), HTTP response status code, and execution latency.

**Retention Policy**: These connection logs are maintained strictly for real-time security monitoring, rate limiting against automated brute-force/DDoS attacks, and infrastructure debugging. All server logs are automatically rotated and permanently purged after 14 days. We never perform cross-site user profiling or link IP addresses to real-world identities.

## 6. Interactive Community Submissions (Discussions, Comments & Reactions)

When engaging with interactive community features:

**Discussions & Q&A (`/discussions`)**: Authenticated members can author new topics, provide technical solutions, post threaded replies, and upvote discussions. Published topics, markdown content, and timestamps are stored in our SQLite database and displayed publicly.

**Article Comments**: You may comment on blog posts either with your verified GitHub profile or via transient anonymous guest submission.

**Post Reactions**: Article reactions (such as Likes or Insightful markers) increment an aggregated numerical counter on the post without storing personal user identifiers.

**Bot Defense Without Surveillance**: We use lightweight, invisible honeypot fields and in-memory rate limiting to filter automated spam bots without employing invasive commercial CAPTCHA widgets.

## 7. Terminal & Shell Interfaces

**In-Browser Terminal (`/terminal`)**: The virtual UNIX web shell runs entirely client-side; command execution history remains in your local browser memory and is discarded upon tab closure.

**Public SSH Access (`ssh daemontalk.com -p 2222`)**: SSH connections run in ephemeral, isolated processes with zero keystroke logging.

## 8. Data Security & Cryptography Standards

We employ defense-in-depth engineering practices to secure all stored and transmitted data:

**Transport Layer Security**: All web traffic is strictly encrypted using TLS 1.3 with Perfect Forward Secrecy (PFS) and HTTP Strict Transport Security (HSTS).

**Session Hashing**: Raw authentication session tokens are never stored in plain text; only salted SHA-256 hashes are persisted in the database.

**Strict Isolation**: The SQLite databases reside on an isolated file system with strict Unix file permissions (`0600`) and prepared parameterized queries to eliminate SQL injection vulnerabilities.

## 9. Data Subject Rights & Self-Service Data Portability

You maintain complete sovereignty over your personal data:

**Self-Service Export**: You can immediately download all your account information, published forum topics, replies, and article comments as a structured JSON file at any time by clicking "Export my data (JSON)" in your profile menu or accessing `/auth/export`.

**Self-Service Account & Data Purge**: You can permanently delete your account at any time via the "Delete account" option in your profile menu. Upon confirmation, your user profile and active sessions are permanently erased from our database, and your public forum discussions and replies are automatically anonymized (`[Deleted User]`) to preserve community knowledge base integrity without linking to your personal identity.

**Manual Requests**: You may also request data correction, export, or removal by contacting us via email at **realdaemontalk@gmail.com** using the address linked to your GitHub account.

## 10. Age Limitations & Children's Privacy

Daemontalk is an engineering and computer science research platform. We do not knowingly collect or solicit personal information from individuals under the age of 13 (or under 16 in certain EU jurisdictions). If we become aware that personal information has been collected from a child without verified parental consent, we will promptly delete that data.

## 11. Server Hosting & International Data Transfers

Our primary infrastructure is self-hosted on secure bare-metal servers located in high-compliance data centers. We do not transfer, route, or mirror personal data across commercial third-party marketing clouds or international data brokers.

## 12. Security Incident Notification Protocols

In the unlikely event of a security incident affecting personal data integrity or confidentiality, Daemontalk will notify affected users and relevant supervisory authorities without undue delay (within 72 hours of becoming aware of the breach) in accordance with applicable data protection laws.

## 13. Policy Amendments & Contact Channels

We may update this Privacy Policy periodically to reflect architectural evolutions, new technical features, or legal requirements. Material updates will be documented on our `/changelog` and indicated by the revision date at the top of this page.

If you have questions, auditing suggestions, or security disclosures regarding our privacy architecture, reach out directly to: **realdaemontalk@gmail.com**.
