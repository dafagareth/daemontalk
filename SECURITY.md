# Security Policy

## Supported Versions

Security updates and patches are exclusively provided for the latest deployment of the `main` branch. Older versions or branches are strictly unsupported.

| Version | Supported |
| ------- | --------- |
| Latest  | Yes       |
| Older   | No        |

## Reporting a Vulnerability

Daemontalk operates under a strict responsible disclosure policy. If you identify a security vulnerability within the codebase, dependencies, or live infrastructure, **do not** create a public GitHub Issue, Pull Request, or publicly broadcast the exploit.

Submit all vulnerability reports privately via email:
**Email:** realdaemontalk@gmail.com

### Reporting Requirements

To ensure a rapid and accurate assessment, your report must include:
1. **Description:** A detailed explanation of the vulnerability and its underlying mechanics.
2. **Reproduction:** Exact step-by-step instructions to reproduce the issue, including endpoint URLs, HTTP request payloads, and required privileges.
3. **Impact:** A technical assessment of the potential exploit impact (e.g., Remote Code Execution, SQL Injection, Cross-Site Scripting, Data Exfiltration).
4. **Environment:** Specify the tools, scripts, or browsers used during discovery.

### Resolution Process

Since Daemontalk is developed and maintained by a solo developer, please allow a reasonable window for triage and remediation. Upon receiving your report, I will follow this protocol:
1. **Triage:** I will do my best to acknowledge receipt of your report within 48-72 hours.
2. **Assessment:** I will personally verify the vulnerability and assess its impact on the system.
3. **Remediation:** A patch will be written, tested, and deployed to the production server as soon as possible.
4. **Disclosure:** Once the issue is resolved, it will be documented in the repository changelog. You will receive full attribution for the responsible disclosure if you wish to be credited.
