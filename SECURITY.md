# Security Policy

## Supported Versions

We actively support the following versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of this project seriously. If you discover a security vulnerability, please report it responsibly.

### How to Report

**Please DO NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **GitHub Security Advisories** (Preferred)
   - Go to the [Security Advisories](../../security/advisories) page
   - Click "Report a vulnerability"
   - Provide detailed information about the vulnerability

2. **Email**
   - Send an email to: `[SECURITY_EMAIL]`
   - Use the subject line: `[SECURITY] Vulnerability Report - [Project Name]`
   - Encrypt sensitive information using our PGP key (available upon request)

### What to Include

Please include the following information in your report:

- **Description**: A clear description of the vulnerability
- **Impact**: The potential impact if exploited
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Affected Versions**: Which versions are affected
- **Suggested Fix**: If you have a suggestion for how to fix it
- **Your Contact Info**: So we can follow up with questions

### What to Expect

1. **Acknowledgment**: We will acknowledge receipt within 48 hours
2. **Initial Assessment**: We will provide an initial assessment within 7 days
3. **Regular Updates**: We will keep you informed of our progress
4. **Resolution**: We aim to resolve critical vulnerabilities within 30 days
5. **Credit**: With your permission, we will credit you in the security advisory

### Security Measures in This Project

This project implements the following security measures:

#### Code Security
- CodeQL static analysis on all PRs and pushes
- Gosec security scanning for Go code
- Dependency vulnerability scanning via govulncheck
- Automated dependency review on pull requests

#### Network Security
- DNS queries use standard, well-documented protocols
- No sensitive data is transmitted or stored
- All DoH/DoT/DoQ connections use TLS 1.3 where supported
- Certificate validation is enforced for encrypted protocols

#### Build Security
- Dependencies are pinned to specific versions
- All builds are reproducible
- Release binaries are checksummed

### Scope

The following are considered in-scope for security reports:

- Vulnerabilities in the application code
- Insecure default configurations
- Dependencies with known vulnerabilities
- Privilege escalation issues
- Information disclosure vulnerabilities

The following are out-of-scope:

- Vulnerabilities in third-party DNS servers being benchmarked
- Issues requiring physical access to the user's machine
- Social engineering attacks
- Denial of service against the benchmarking tool itself

### Safe Harbor

We support safe harbor for security researchers who:

- Make a good faith effort to avoid privacy violations
- Avoid disruption to production systems
- Provide us reasonable time to address the issue before disclosure
- Do not exploit vulnerabilities beyond what's necessary to demonstrate them

We will not pursue legal action against researchers who follow these guidelines.

## Security Best Practices for Users

When using this tool:

1. **Keep Updated**: Always use the latest version
2. **Verify Downloads**: Check checksums of downloaded binaries
3. **Network Awareness**: Be aware that DNS queries will be sent to configured servers
4. **Import Trusted Sources**: Only import server lists from trusted sources

## Security Updates

Security updates will be released as:

- Patch releases for non-breaking fixes
- GitHub Security Advisories with CVE identifiers when applicable
- Announcements in the project's release notes

---

*This security policy is subject to change. Last updated: December 2025*
