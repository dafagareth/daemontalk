# Terms of Use

**Last updated: September 5, 2026** · **Version: 2.3**

Welcome to Daemontalk. By accessing our website, reading published technical articles, engaging in discussions on Socket (`/socket`), leaving comments, connecting to our public SSH TUI gateway (`ssh daemontalk.com -p 2222`), requesting content via command-line utilities (`curl`), or subscribing to syndication feeds (RSS/JSON Feed) ("Services"), you agree to be bound by these Terms of Use and our [Privacy Policy](/privacy). If you do not agree to these terms, please discontinue using the platform immediately.

---

## 1. Acceptance of Terms & Scope of Services

Daemontalk is an independent systems and software engineering publication, open tech notebook, and developer community platform focusing on operating systems, backend architectures, Linux kernel internals, and engineering culture. These Services include:

- The public web interface (`daemontalk.com` and `/id` paths), encompassing all technical articles, editorial essays, and architectural deep-dives.
- The **Socket** Developer Community & Discussion Forum (`/socket`).
- Article comment threads authenticated via official GitHub OAuth or guest aliases.
- The interactive terminal reader interface accessible via our **Public SSH TUI Gateway** (`ssh daemontalk.com -p 2222`).
- Command-line friendly text-streaming endpoints (such as `curl daemontalk.com/daily`, `/recipes`, `/p/:slug`).
- Public syndication feeds (RSS 2.0 at `/rss.xml` and JSON Feed v1.1 at `/feed.json`).

*(Architectural note: The legacy in-browser web terminal `/terminal` has been officially deprecated and completely removed from the platform to maintain a lightweight, high-performance Go SSR binary, tight security boundaries, and an authentic terminal experience via SSH and curl).*

## 2. Intellectual Property & Content Licensing

Intellectual property rights and licensing across Daemontalk are transparently defined to foster an open, educational ecosystem:

**Technical Articles & Editorial Essays**: All original technical articles, system analysis dispatches, and educational essays published under `content/posts/` are licensed under the **[Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International (CC BY-NC-SA 4.0)](https://creativecommons.org/licenses/by-nc-sa/4.0/)** license. You are free to share, copy, and adapt the text for non-commercial purposes, provided clear attribution (Dafa Gareth / Daemontalk) and a direct backlink to the original article are included.

**Code Snippets & Configuration Blueprints**: All source code snippets, Linux kernel sysctl parameters, shell scripts, and build configurations embedded within technical articles are provided under permissive open-source terms (MIT / Unlicense). You are free to copy, modify, and integrate them into personal or commercial production systems without royalties.

**Platform Source Code & Server Binaries**: The underlying software architecture of Daemontalk (Go backend, compiled Templ views, Tailwind stylesheets, Wish SSH daemon, and Bubble Tea TUI) is licensed under the **[PolyForm Noncommercial License 1.0.0](https://polyformproject.org/licenses/noncommercial/1.0.0)**. You may review, audit, fork, and adapt the code for personal study, academic research, or non-commercial deployments. You may not commercialize, resell, or distribute the platform as a commercial SaaS product without prior written consent.

**Proprietary Branding & Visual Assets**: The **Daemontalk** name, `daemontalk.com` domain, mechanical daemon mascot, design tokens, layout geometry, and visual brand identity are proprietary intellectual property (Copyright © 2026 Dafa Gareth. *All Rights Reserved*). Cloning or utilizing these brand assets for confusing or fraudulent purposes is strictly prohibited.

## 3. User-Generated Content, Comments & Socket Forum

**License Grant on Contributions**: When you submit questions, architectural solutions, technical code samples, or replies within the Socket forum (`/socket`) or article comment sections, you grant Daemontalk a perpetual, worldwide, non-exclusive, royalty-free license to display, index, format, and distribute your contributions as part of the public tech knowledge base.

**Originality & Legal Responsibility**: You retain copyright ownership over your original submissions. You represent and warrant that your contributions are your own work or that you possess the necessary rights and permissions to post them, without infringing third-party intellectual property, non-disclosure agreements (NDAs), or proprietary employer code.

**Preservation & Account Deletion Anonymization**: If you choose to delete your account, your personal profile data is permanently purged from our servers. To maintain the educational coherence and continuity of community problem-solving threads, previously submitted questions and answers will remain archived in anonymized form (`[Deleted User]`).

## 4. User Accounts & GitHub OAuth Authentication

**Developer-Centric Authentication**: Membership registration on Daemontalk is handled exclusively via official GitHub OAuth. We request only minimal, public profile identifiers (GitHub ID, username, avatar URL, and primary verified email) for developer identification. We never request, access, or store your GitHub account passwords.

**Account Security & Session Activity**: You are solely responsible for maintaining the confidentiality of your credentials and for all activities conducted under your authenticated session. If you suspect unauthorized access, you can instantly revoke Daemontalk's OAuth access via your GitHub security settings.

**Account Non-Transferability**: User accounts are personal, non-transferable, and may not be bought, sold, or shared with third parties.

## 5. Community Code of Conduct & Prohibited Usage

When participating in Socket discussions, creating topics, or submitting comments:

**Constructive Technical Discourse**: Maintain professional, evidence-based, civil, and constructive dialogue. Focus on technical problem-solving, reproducible benchmarks, and architectural trade-offs. Challenge technical concepts and code objectively, never individuals.

**Prohibited Submissions & Spam**: You may not publish commercial advertisements, unsolicited affiliate links, automated bot campaigns, low-quality AI-generated spam, malicious code exploits, or phishing URLs.

**Anti-Sybil & Voting Integrity**: Fabricating duplicate accounts to manipulate topic upvotes, falsely mark solutions as resolved, artificially inflate forum standing, or circumvent rate-limiting mechanisms is strictly forbidden.

**Harassment & Privacy Protection**: Doxxing personal information, harassment, stalking, bullying, or hate speech will result in immediate and permanent account termination.

## 6. System Integrity, SSH Gateway & Command-Line (CLI) Access

Daemontalk provides direct terminal access to deliver an authentic, distraction-free technical reading experience:

**Public SSH TUI Gateway (`ssh daemontalk.com -p 2222`)**: The SSH TUI reader is provided for interactive article and briefing browsing directly from your local terminal. You agree to use SSH sessions reasonably. Attempting to escape the TUI sandbox (pty breakout), injecting malicious terminal escape sequences, inducing server panics, or utilizing the SSH port as an open proxy or tunnel is strictly prohibited.

**Command-Line CLI Endpoints (`curl`)**: Plain-text streaming endpoints (`/daily`, `/recipes`, `/p/:slug`) are offered for console convenience. Automated access is permitted provided it uses polite User-Agents and adheres to reasonable rate limits without causing infrastructure degradation.

**Infrastructure Security Prohibitions**: You must not conduct unauthorized vulnerability scans, distributed denial-of-service (DDoS/DoS) attacks, brute-force attacks against SSH or HTTP ports, or aggressive web scraping that threatens platform stability for other users.

## 7. Technical Disclaimer & "As-Is" Terms

All technical articles, system administration tutorials, Linux kernel sysctl parameters, shell scripts, database configurations, and performance benchmarks published on Daemontalk are provided purely for educational and informational purposes.

Software environments and operating systems evolve rapidly. A configuration that works reliably in our test environment may carry unforeseen side effects in your unique architecture. You assume full and sole responsibility for auditing, testing, and benchmarking any command, configuration, or code sample in an isolated staging or sandbox environment before deploying it to production infrastructure.

## 8. Limitation of Liability

To the maximum extent permitted by applicable law, Daemontalk (operated by Dafa Gareth) and its article contributors shall not be liable for any direct, indirect, incidental, special, exemplary, or consequential damages.

This includes, without limitation, hardware failure, kernel panics, system boot failures, data corruption or loss, production downtime, loss of business revenue, or security breaches arising from or in connection with your application of materials, scripts, or advice found on this platform.

## 9. Content Moderation & Account Revocation

**Maintainer Authority**: We reserve the right to review, edit, lock, unlist, or remove any forum topic, reply, or comment that violates these terms, degrades community health, or breaches legal regulations, without prior notice.

**Account Revocation**: We reserve the right to suspend or permanently ban user accounts that engage in persistent abuse, spamming, or security attacks against the platform.

**Self-Service Account Termination**: You retain complete autonomy over your account and may permanently delete your profile at any time directly through the "Delete account" option in your user profile dropdown.

## 10. Third-Party Links & External References

Articles and user submissions frequently reference third-party resources, such as GitHub repositories, RFC standards, academic papers, and official open-source documentation. Daemontalk does not control, endorse, or assume liability for the accuracy, uptime, or privacy practices of external websites.

## 11. Service Availability & Platform Evolution

We reserve the right to modify, suspend, or discontinue any feature, endpoint, or service component (including the Socket forum, SSH gateway, or CLI endpoints) at any time without liability, as part of continuous server maintenance, security patching, and architectural enhancements.

## 12. Governing Law, Dispute Resolution & Official Contact

These Terms shall be governed by and construed in accordance with the laws of the Republic of Indonesia. Any dispute arising out of or related to these terms shall be addressed through amicable, good-faith technical consultation.

For inquiries, copyright notices, or commercial licensing requests, contact us directly at:
- **Official Email**: [realdaemontalk@gmail.com](mailto:realdaemontalk@gmail.com)
- **Project Repository & Discussions**: [github.com/dafagareth/daemontalk](https://github.com/dafagareth/daemontalk)
