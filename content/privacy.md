# Privacy Policy

**Last updated: August 27, 2026**

Daemontalk operates under a strict zero-tracker philosophy. We believe technical publications should prioritize reader autonomy, transparency, and computational minimalism. We do not track, profile, monetize, or harvest personal data.

## 1. Zero Tracking & Self-Hosted Infrastructure

- **No Behavioral Trackers**: This site operates without Google Analytics, Meta/Facebook Pixels, tracking beacons, advertising cookies, or browser fingerprinting scripts.
- **No Third-Party Ad Networks**: We do not display commercial advertisements, sponsor banners, or affiliate trackers, and we never sell user data to data brokers.
- **Self-Hosted Assets**: All stylesheets, scripts, and media are served directly from our dedicated Go binary server, eliminating external tracking surfaces.

## 2. Client-Side Storage & Local Preferences

We utilize standard client-side browser storage mechanisms (`localStorage` and `sessionStorage`) exclusively on your local machine to preserve your personal viewing preferences:

- **Saved Dispatches (`bookmarks`)**: The list of bookmarked posts stored in your local reading ledger.
- **Theme Selection (`theme`)**: Your active color scheme preference (*Light, Dark, or Sepia mode*).
- **Typography & Accessibility**: Font size scaling, font family selection (*Serif or Sans*), and warm screen tint intensity.
- **Visited Dispatches Indicator (`readPosts`)**: A local array of recently visited article slugs (capped at 200 items) used solely to indicate previously read content on index pages.
- **Session Animations (`visited`)**: A temporary session token in `sessionStorage` to prevent redundant entry animations on subsequent page navigations.

**Data Sovereignty**: None of this client-side preference data is ever synced, uploaded, or transmitted to our backend database or any third-party entity. You may purge this data at any time by clearing your browser's site storage.

## 3. Ephemeral Server Logs & Telemetry

When connecting to Daemontalk over HTTP or SSH, our self-hosted server records standard connection metadata:

- Connecting IP address.
- Request path, HTTP method, and timestamp.
- User-Agent string and Referrer header (if transmitted by your client).
- HTTP response status code and execution latency.

**Retention Policy**: These connection logs are maintained strictly for real-time security monitoring, rate limiting against automated brute-force/DDoS attacks, and infrastructure debugging. All server logs are automatically rotated and permanently purged after **14 days**. We never perform cross-site user profiling or link IP addresses to real-world identities.

## 4. Interactive Submissions (Comments, Reactions & Guestbook)

When engaging with interactive community features:

- **Public Discussions**: When you submit a comment or guestbook entry, the author handle (Callsign) and message body you supply are stored in our local SQLite database and published publicly.
- **Post Reactions**: Article reactions (such as Likes or Insightful markers) increment an aggregated numerical counter on the post without storing personal user identifiers.
- **No Compulsory Registration**: You are never required to provide passwords, email addresses, phone numbers, or OAuth social accounts to participate.
- **Bot Defense Without Surveillance**: We use lightweight, invisible honeypot fields to filter automated spam bots without employing invasive commercial CAPTCHA widgets.

## 5. Terminal & Shell Interfaces

- **In-Browser Terminal (`/terminal`)**: The virtual UNIX web shell runs entirely client-side; command execution history remains in your local browser memory and is discarded upon tab closure.
- **Public SSH Access (`ssh daemontalk.com -p 2222`)**: SSH connections run in ephemeral, isolated processes with zero keystroke logging.

## 6. Data Subject Rights & Content Removal

You maintain the right to request the correction or permanent removal of any comment or guestbook transmission you have submitted. To request content removal, contact us at **realdaemontalk@gmail.com** specifying the article slug and your submitted author handle.

## 7. External References

Our technical dispatches frequently cite open-source repositories (GitHub), documentation hubs, and academic papers (arXiv, IEEE). When following external hyperlinks, your browsing activity is governed by the privacy policies of the destination services.

## 8. Inquiries & Contact

If you have questions, auditing suggestions, or security disclosures regarding our privacy architecture, reach out directly to: **realdaemontalk@gmail.com**.
