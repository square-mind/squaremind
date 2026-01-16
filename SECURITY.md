# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

We take security seriously at Squaremind. If you discover a security vulnerability, please report it responsibly.

### How to Report

1. **Do NOT open a public GitHub issue**
2. Email security@squaremind.xyz with details
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Any suggested fixes

### What to Expect

- **Acknowledgment**: Within 48 hours
- **Initial Assessment**: Within 1 week
- **Resolution Timeline**: Depends on severity
  - Critical: 24-48 hours
  - High: 1 week
  - Medium: 2 weeks
  - Low: Next release

### Scope

The following are in scope:
- Cryptographic identity (Ed25519) implementation
- Gossip protocol message handling
- Consensus mechanism
- Task market bidding logic
- Reputation calculation
- CLI command injection
- SDK vulnerabilities

### Out of Scope

- Vulnerabilities in dependencies (report to upstream)
- Social engineering attacks
- Physical attacks
- Denial of service (unless protocol-level)

## Security Best Practices

### For Users

1. **Protect API Keys**: Never commit API keys to version control
2. **Use Environment Variables**: Store secrets in environment variables
3. **Validate Agent Identity**: Verify SID signatures in production
4. **Monitor Reputation**: Watch for anomalous reputation changes
5. **Network Security**: Use TLS for gossip protocol in production

### For Developers

1. **Input Validation**: Validate all external input
2. **Cryptographic Operations**: Use standard library crypto
3. **Error Handling**: Don't leak sensitive info in errors
4. **Dependency Updates**: Keep dependencies current
5. **Code Review**: All changes require review

## Cryptographic Details

Squaremind uses the following cryptographic primitives:

| Component | Algorithm | Library |
|-----------|-----------|---------|
| Identity | Ed25519 | golang.org/x/crypto/ed25519 |
| Hashing | SHA-256 | crypto/sha256 |
| Random | crypto/rand | Go standard library |

## Acknowledgments

We appreciate responsible disclosure. Security researchers who report valid vulnerabilities will be acknowledged (with permission) in our release notes.
